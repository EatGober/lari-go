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
	"time"

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

			ctr := 0

			// Parse the date-time string
			t, err := time.Parse(time.RFC3339, data.ScheduleDateTimeString)
			if err != nil {
				fmt.Println("Error parsing date-time:", err)
				return
			}

			// Format the time as MM/DD/YY HH:MM AM/PM
			formattedTime := t.Format("01/02/06 03:04 PM")

			// Print formatted time
			fmt.Println("Formatted Time:", formattedTime)

			for _, appt := range apptList {
				ptURL := "http://us-east.performave.com:3001/confirm/" + (apptId) + "/" + strconv.Itoa(appt.PatientID)

				text := "Hello! A waitlist spot has opened up on " + formattedTime + "! Click the link to take the spot: " + ptURL

				if ctr == 0 {
					sms.SendMessage(appt.PatientPhone, text)
					ctr++
				}
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
			context.Redirect(http.StatusFound, os.Getenv("ADDR_SUCCESS"))
		} else {
			// redirect to fail page
			context.Redirect(http.StatusFound, os.Getenv("ADDR_FAIL"))
		}

	})

	// Router link /confirm/USERHASH -> confirm user, delete other users, fwd user to conf page

	router.Run(":3001")
}
