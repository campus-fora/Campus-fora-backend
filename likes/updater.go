package likes

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
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
		false,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Error in registering updater consumer: %s", err)
		return
	}

	updaterDB, err := openDBConn()

	var likeCountBuffer LikeCountBuffer
	likeCountBuffer.hmap = make(map[uuid.UUID]LikeCountBufferValues)
	var buffer []amqp.Delivery

	processingInterval := time.Duration(viper.GetInt("MQ.PROCESSING_INTERVAL")) * time.Second
	timer := time.NewTicker(processingInterval)

	voteRequestsCh := make(chan newVoteRequest)
	createWorkers(voteRequestsCh, &likeCountBuffer)
	// for msg := range msgs {
	// 	buffer = append(buffer, msg)
	// 	var voteRequest newVoteRequest
	// 	err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&voteRequest)
	// 	if err != nil {
	// 		log.Println("Error in decoding vote request struct")
	// 		continue
	// 	}
	// 	voteRequestsCh <- voteRequest
	// 	log.Print("Recieved vote request: ", voteRequest.PostID, voteRequest.UserID, voteRequest.VoteType)
	// 	// err = updateUserLike(voteRequest.PostID, voteRequest.UserID, int(voteRequest.VoteType))
	// 	// if err != nil {
	// 	// 	log.Fatal("Error updating user like: ", err)
	// 	// }
	// }

	defer processBufferedMessages(buffer, &likeCountBuffer, updaterDB)
	defer acknowledgeMessages(buffer)
	for {
		select {
			case message := <-msgs:
				buffer = append(buffer, message)
				var voteRequest newVoteRequest
				err := gob.NewDecoder(bytes.NewReader(message.Body)).Decode(&voteRequest)
				if err != nil {
					log.Println("Error in decoding vote request struct")
					continue
				}
				voteRequestsCh <- voteRequest
				log.Print("Recieved vote request: ", voteRequest.PostID, voteRequest.UserID, voteRequest.VoteType)
			case <-timer.C:
				processBufferedMessages(buffer, &likeCountBuffer, updaterDB)
				acknowledgeMessages(buffer)
				buffer = nil
		}
	}
}

func createWorkers(voteRequestsCh chan newVoteRequest, likeCountBuffer *LikeCountBuffer) {
	workerCount := viper.GetInt("MQ.WORKER_COUNT")
	for i := 0; i < workerCount; i++ {
		go worker(voteRequestsCh, i+1, likeCountBuffer)
	}
}

func worker(voteRequestsCh chan newVoteRequest, id int, likeCountBuffer *LikeCountBuffer) {
	workerDB, err := openDBConn()
	if err != nil {
		log.Printf("Error in opening db connection for worker %d: %s", id, err)
		return
	}

	for voteRequest := range voteRequestsCh {
		log.Println("worker", id, ":", voteRequest)
		err, prevVote, latestTimeStamp := updateUserLikeStatus(workerDB, voteRequest.PostID, voteRequest.UserID, voteRequest.VoteType, voteRequest.LatestReqTime)
		if(!latestTimeStamp) {
			continue
		}
		likeCountBuffer.updatelikeCountBuffer(voteRequest, workerDB, prevVote)
		if err != nil {
			log.Printf("WORKER %d : Error in updating like status for user %d and post %d", id, voteRequest.UserID, voteRequest.PostID)
		}
	}
}
