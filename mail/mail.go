package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Mail struct {
	To      []string
	Subject string
	Body    string
}

func (mail *Mail) BuildMessage() []byte {
	message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message += fmt.Sprintf("From: rggrrffrtg<%s>\r\n", sender)
	message += fmt.Sprintf("Subject: %s | efffeefefefe\r\n", mail.Subject)

	// If mass mailing, BCC all the users
	if len(mail.To) == 1 {
		message += fmt.Sprintf("To: %s\r\n\r\n", mail.To[0])
	} else {
		message += "To: Undisclosed Recipients\r\n\r\n"
	}

	message += strings.Replace(mail.Body, "\n", "<br>", -1)
	message += "<br><br>--<br>rrggrgrge<br>"
	message += "feffeffefefeeef<br><br>"
	message += "<small>This is an auto-generated email. Please do not reply.</small>"

	return []byte(message)
}

func Service(mailQueue chan Mail) {
	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", user, pass, host)

	for mail := range mailQueue {
		message := mail.BuildMessage()
		to := mail.To

		if err := smtp.SendMail(addr, auth, sender, to, message); err != nil {
			fmt.Printf("Error sending mail: %s", to)
			fmt.Printf("Error: %s", err)
		}
	}
}
