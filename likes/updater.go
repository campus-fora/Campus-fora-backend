package likes

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type PostLikeCountBuffer struct {
	mu sync.Mutex
	hmap map[int]int
}

var postLikeCountBuffer PostLikeCountBuffer

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
	// open db connection
	for voteRequest := range voteRequestsCh {
		log.Println("worker",id, ":" , voteRequest)
		time.Sleep(10)
	}
}

func batchUpdater() {
	defer ch.BatchUpdater.Close()
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
	q, err := ch.BatchUpdater.QueueDeclare(
		viper.GetString("MQ.BATCHUPDATERQUEUE"), // name
		true,                                    // durable
		false,                                   // delete when unused
		false,                                   // exclusive
		false,                                   // no-wait
		nil,                                     // arguments
	)
	if err != nil {
		log.Printf("Error in declaring updater queue: %s", err)
		return
	}
	err = ch.BatchUpdater.QueueBind(
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
	msgs, err := ch.BatchUpdater.Consume(
		q.Name,          // queue
		"batch_updater", // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Printf("Error in registering updater consumer: %s", err)
		return
	}

	processingInterval := time.Duration(viper.GetInt("MQ.PROCESSING_INTERVAL")) * time.Second
	timer := time.NewTicker(processingInterval)
	var buffer []amqp.Delivery

	for {
		select {
		case message := <-msgs:
			updatePostLikeCountBuffer(message)
		case <-timer.C:
			// postLikeCountBuffer.mu.Lock()
			processBufferedMessages(buffer)
			acknowledgeMessages(buffer)
			buffer = nil
			postLikeCountBuffer.hmap = nil
			// postLikeCountBuffer.mu.Unlock()
		}
	}
}

func updatePostLikeCountBuffer(message amqp.Delivery) {
	var voteRequest newVoteRequest
	err := gob.NewDecoder(bytes.NewReader(message.Body)).Decode(&voteRequest)
	if err != nil {
		log.Println("Error in decoding vote request struct")
		return
	}
	// postLikeCountBuffer.mu.Lock()
	// defer postLikeCountBuffer.mu.Unlock()
	// postLikeCountBuffer.hmap[int(voteRequest.PostID)] += int(voteRequest.VoteType)
}

func processBufferedMessages(buffer []amqp.Delivery) {
	if len(buffer) == 0 {
		return
	}
	log.Println("Processing batch:")
	for _, msg := range buffer {
		var voteRequest newVoteRequest
		err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&voteRequest)
		if err != nil {
			log.Println("Error in decoding vote request struct")
			continue
		}
		log.Print("Message: ", voteRequest)
		// updateBatchLikeCount(voteRequest.)
	}
	log.Println()
}

func acknowledgeMessages(buffer []amqp.Delivery) {
	for _, msg := range buffer {
		msg.Ack(false)
	}
}
