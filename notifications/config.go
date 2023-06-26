package notifications

import (
	"context"
	"log"
	_ "github.com/campus-fora/config"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Notification struct {
	title  string
	body   string
	link   string
	tokens []string
}

type fcmMessaging struct {
	fcm_client *messaging.Client
	notif_chan chan Notification
}

var Notification_client fcmMessaging

func init() {
	notif_chan := make(chan Notification)
	Notification_client.notif_chan = notif_chan

	opt := option.WithCredentialsFile("firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	messaging, err := app.Messaging(context.TODO())
	if err != nil {
		log.Fatalf("messaging: %s", err)
	}
	Notification_client.fcm_client = messaging
	log.Println("Notification client initialized")
	go sender()
}
