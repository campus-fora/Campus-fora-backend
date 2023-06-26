package notifications

import "firebase.google.com/go/messaging"

func GeneratePushNotification(title string, body string, link string, tokens []string) *messaging.MulticastMessage {
	return &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Webpush: &messaging.WebpushConfig{
			FcmOptions: &messaging.WebpushFcmOptions{
				Link: link,
			},
			Notification: &messaging.WebpushNotification{
				Icon: "",
			},
		},
		Tokens: tokens,
	}
}
