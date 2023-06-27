package likes

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/spf13/viper"
)

func updater() {
	defer ch.Updater.Close()
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
	q, err := ch.Updater.QueueDeclare(
		viper.GetString("MQ.UPDATERQUEUE"), // name
		true,                               // durable
		false,                              // delete when unused
		false,                              // exclusive
		false,                              // no-wait
		nil,                                // arguments
	)
	if err != nil {
		log.Printf("Error in declaring updater queue: %s", err)
		return
	}
	err = ch.Updater.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error in binding updater queue: %s", err)
		return
	}
	msgs, err := ch.Updater.Consume(
		q.Name,    // queue
		"updater", // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Error in registering updater consumer: %s", err)
		return
	}

	voteRequestsCh := make(chan newVoteRequest)
	createWorkers(voteRequestsCh)
	for msg := range msgs {
		var voteRequest newVoteRequest
		err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&voteRequest)
		if err != nil {
			log.Println("Error in decoding vote request struct")
			continue
		}
		voteRequestsCh <- voteRequest
		log.Print("Recieved vote request: ", voteRequest.PostID, voteRequest.UserID, voteRequest.VoteType)
		// err = updateUserLike(voteRequest.PostID, voteRequest.UserID, int(voteRequest.VoteType))
		// if err != nil {
		// 	log.Fatal("Error updating user like: ", err)
		// }
	}
}

func createWorkers(voteRequestsCh chan newVoteRequest) {
	workerCount := viper.GetInt("MQ.WORKER_COUNT")
	for i := 0; i < workerCount; i++ {
		go worker(voteRequestsCh, i+1)
	}
}

func worker(voteRequestsCh chan newVoteRequest, id int) {
	workerDB, err := openDBConn()
	if err != nil {
		log.Printf("Error in opening db connection for worker %d: %s", id, err)
		return
	}
	log.Print("Worker ", id, workerDB)

	for voteRequest := range voteRequestsCh {
		log.Println("worker", id, ":", voteRequest)
		err = updateUserLikeStatus(workerDB, voteRequest.PostID, voteRequest.UserID, voteRequest.VoteType)
		if err != nil {
			log.Printf("WORKER %d : Error in updating like status for user %d and post %d", id, voteRequest.UserID, voteRequest.PostID)
		}
	}
}
