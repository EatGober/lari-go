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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RunServer() {
	router := gin.Default()

	router.PUT("/update", func(context *gin.Context) {
		status := context.GetHeader("Status")
		var data domain.CancelledAppointment

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

			// context.Params

			addr := os.Getenv("MDW_ADDR") + "/waitlist/195900/1/71"

			fmt.Println(addr)
			resp, err := http.Get(addr)

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

			apptId := scheduler.CreateSlot(apptList)

			for _, appt := range apptList {
				ptURL := "http://us-east.performave.com:3001/" + (apptId) + "/" + strconv.Itoa(appt.PatientID)

				text := "Hello! A waitlist spot has opened up on " + data.ScheduleDateTimeString + "! Click the link to take the spot: " + ptURL

				// sms.SendMessage(appt.PatientPhone, text)
				sms.DummyMessage(appt.PatientPhone, text)
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
