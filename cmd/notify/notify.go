package notify

import (
	// "lari-go/messaging"

	"net/http"
	// "encoding/json"
	// "log"

	"github.com/gin-gonic/gin"
	twiml "github.com/twilio/twilio-go/twiml"
)

func RunServer() {
	router := gin.Default()

	router.POST("/sms", func(context *gin.Context) {
		message := twiml.MessagingMessage{
			Body: "Thanks! Your appointment has been scheduled for",
		}

		twimlResult, err := twiml.Messages([]twiml.Element{message})

		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
		} else {
			context.Header("Content-Type", "text/xml")
			context.String(http.StatusOK, twimlResult)
		}
	})

	router.Run(":3000")
}
