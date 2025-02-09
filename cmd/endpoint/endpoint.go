package endpoint

import (
	"fmt"
	"io"
	"lari-go/internal/domain"
	"lari-go/internal/scheduler"
	"lari-go/internal/sms"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RunServer() {
	router := gin.Default()

	router.PUT("/update", func(context *gin.Context) {
		status := context.GetHeader("Status")
		var data domain.Appointment

		if status == "" {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Status header is required",
			})
		} else if status == "cancelled" {

			if err := context.ShouldBindJSON(&data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				// return
			}

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error getting current directory: %v", err)
			}

			if len(cwd) == 0 {
				cwd = "root/lari-go"
			}

			// Construct the full path to .env
			envPath := cwd + "/.env"

			fmt.Println("Loading environment variables...")
			err = godotenv.Load(envPath)

			if err != nil {
				log.Fatal("Error loading .env")
				// return
			}

			fmt.Println(os.Getenv("MDW_ADDR") + "/waitlist/195900")

			resp, err := http.Get(os.Getenv("MDW_ADDR") + "/waitlist/195900")

			if err != nil {
				log.Fatal("Error getting waitlist data")
				// return
			}

			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				// return
			}

			bodyString := string(bodyBytes)

			apptList := domain.ConstructApptLists(bodyString)

			scheduler.CreateSlot(apptList)

			for _, appt := range apptList {
				// appt.PatientPhone
				sms.DummyMessage(fmt.Sprint(appt.PatientPhone), "YAY APPOINTMENT?!?!?!?!")
			}

		} else {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Status header value is invalid",
			})
		}

	})

	router.GET("/confirm/:timeslotid/:patientid", func(context *gin.Context) {
		patientId := context.Param("patientid")
		timeslotId := context.Param("timeslotid")

		success := scheduler.Validate(patientId, timeslotId)

		if success {
			// redirect to success page
			scheduler.Remove(timeslotId)
			context.Redirect(http.StatusFound, "https://google.com")
		} else {
			// redirect to fail page
			context.Redirect(http.StatusFound, "https://bing.com")
		}

	})

	// Router link /confirm/USERHASH -> confirm user, delete other users, fwd user to conf page

	router.Run(":3001")
}
