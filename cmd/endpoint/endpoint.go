package notify

import (
	// "lari-go/messaging"

	// "encoding/json"
	// "log"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	router := gin.Default()

	router.POST("/sms", func(context *gin.Context) {

	})

	router.POST("/confirm/:id", func(context *gin.Context) {

	})

	// Router link /confirm/USERHASH -> confirm user, delete other users, fwd user to conf page

	router.Run(":3000")
}
