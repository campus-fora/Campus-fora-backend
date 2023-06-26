package notifications

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func sender() {
	for notification := range Notification_client.notif_chan {
		// log.Print("notification ->",notification)
		_, err := Notification_client.fcm_client.SendMulticast(context.Background(),
			GeneratePushNotification(
				notification.title,
				notification.body,
				notification.link,
				notification.tokens,
			))
		if err != nil {
			log.Fatal("error encountered by notification sender: ", err)
		}
	}
	log.Println("Notification sender stopped")
}

func TestNotification(ctx *gin.Context) {
	var notificationRequest Notification
	log.Println("notificationRequest ->", ctx.Request.Body)
	if err := ctx.ShouldBindJSON(&notificationRequest); err != nil {
		log.Print("error ->",err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notificationRequest = Notification{
		title:  "Test Notification",
		body:   "This is a test notification",
		link:   "https://www.google.com",
		tokens: []string{"f4IFaI5cbOL4nDrnVXihvU:APA91bH8VDWTk59da9iQZWPyiThrNMyOmk9yBB5-4ZyKtyRcu_a0CAOl0sWHzYwaHLDqf35vUVh-ag_0l8WTbVwaS9uIIDp4iOMvM0b6QLD7VeX90Epjh3C2pEJeba02yg8cXkFhg1qv"},
	}
	log.Println("notificationRequest ->", notificationRequest)
	Notification_client.notif_chan <- notificationRequest
}
