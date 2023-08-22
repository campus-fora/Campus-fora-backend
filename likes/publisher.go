package likes

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func publisher() {
	defer ch.Publisher.Close()
	exchangeName := viper.GetString("MQ.EXCHANGE")
	err := ch.Publisher.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Error in declaring exchange: %s", err)
		return
	}

	for voteRequest := range voteCh {
		var body bytes.Buffer
		err := gob.NewEncoder(&body).Encode(voteRequest)
		if err != nil {
			log.Println("Error in encoding vote request struct")
			continue
		}
		publisherTimeout := time.Duration(viper.GetInt("MQ.PUBLISHERTIMEOUT")) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), publisherTimeout)
		defer cancel()
		err = ch.Publisher.PublishWithContext(ctx,
			exchangeName, // exchange
			"",           // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body.Bytes(),
			})
		if err != nil {
			log.Println("Error in publishing vote request")
		}
	}
}
