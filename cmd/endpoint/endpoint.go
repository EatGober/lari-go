package endpoint

import (
	"fmt"
	"lari-go/internal/scheduler"
	"net/http"

	// "encoding/json"
	// "log"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	fmt.Println("text")
	router := gin.Default()

	// use router.POST when actually implementing
	router.GET("/sms", func(context *gin.Context) {

	})

	router.GET("/confirm/:timeslotid/:patientid", func(context *gin.Context) {
		patientId := context.Param("patientid")
		timeslotId := context.Param("timeslotid")

		success := scheduler.Validate(patientId, timeslotId)

		if success {
			// redirect to success page
			context.Redirect(http.StatusFound, "https://google.com")
		} else {
			// redirect to fail page
			context.Redirect(http.StatusFound, "https://bing.com")
		}

	})

	// Router link /confirm/USERHASH -> confirm user, delete other users, fwd user to conf page

	router.Run(":3000")
}
