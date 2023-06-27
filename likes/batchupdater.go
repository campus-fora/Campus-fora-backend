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

type LikeCountBuffer struct {
	mu sync.Mutex
	hmap map[uint]int
}

func (likeCountBuffer *LikeCountBuffer) updatelikeCountBuffer(message amqp.Delivery) {
	var voteRequest newVoteRequest
	err := gob.NewDecoder(bytes.NewReader(message.Body)).Decode(&voteRequest)
	if err != nil {
		log.Println("Error in decoding vote request struct")
		return
	}
	likeCountBuffer.mu.Lock()
	defer likeCountBuffer.mu.Unlock()
	likeCountBuffer.hmap[voteRequest.PostID] += int(voteRequest.VoteType)
}

func openBatchUpdaterQueue() (<-chan amqp.Delivery, error) {
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return msgs, nil
}

func batchUpdater() {
	defer ch.BatchUpdater.Close()
	processingInterval := time.Duration(viper.GetInt("MQ.PROCESSING_INTERVAL")) * time.Second
	timer := time.NewTicker(processingInterval)

	var likeCountBuffer LikeCountBuffer
	likeCountBuffer.hmap = make(map[uint]int)
	var buffer []amqp.Delivery

	msgs, err := openBatchUpdaterQueue()
	if err != nil {
		return
	}

	defer processBufferedMessages(buffer,&likeCountBuffer)
	defer acknowledgeMessages(buffer)
	for {
		select {
		case message := <-msgs:
			likeCountBuffer.updatelikeCountBuffer(message)
			buffer = append(buffer, message)
		case <-timer.C:
			processBufferedMessages(buffer,&likeCountBuffer)
			acknowledgeMessages(buffer)
			buffer = nil
		}
	}
}

func processBufferedMessages(buffer []amqp.Delivery, likeCountBuffer *LikeCountBuffer) {
	if len(buffer) == 0 {
		return
	}
	logBuffer(buffer)

	likeCountBuffer.mu.Lock()
	defer likeCountBuffer.mu.Unlock()
	for pid, likeCnt := range likeCountBuffer.hmap {
		updateBatchLikeCount(pid, likeCnt)
	}
	likeCountBuffer.hmap = make(map[uint]int)
}

func acknowledgeMessages(buffer []amqp.Delivery) {
	if len(buffer) == 0 {
		return
	}
	for _, msg := range buffer {
		msg.Ack(false)
	}
}

func logBuffer(buffer []amqp.Delivery) {
	log.Println("Processing batch:")
	for _, msg := range buffer {
		var voteRequest newVoteRequest
		err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&voteRequest)
		if err != nil {
			log.Println("Error in decoding vote request struct")
			continue
		}
		log.Print("Message: ", voteRequest)
	} 
}