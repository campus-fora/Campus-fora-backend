package likes

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type LikeCountBufferValues struct{
	likeCnt   int
	dislikeCnt int
}
type LikeCountBuffer struct {
	mu   sync.Mutex
	hmap map[uuid.UUID] LikeCountBufferValues
}

func (likeCountBuffer *LikeCountBuffer) updatelikeCountBuffer(voteRequest newVoteRequest, db *gorm.DB, prevVote int) {
	// var voteRequest newVoteRequest
	// err := gob.NewDecoder(bytes.NewReader(message.Body)).Decode(&voteRequest)
	// if err != nil {
	// 	log.Println("Error in decoding vote request struct")
	// 	return
	// }
	// oldVoteStatus, err:= fetchLikeStatusWithoutContext(db,uint(voteRequest.PostID), uint(voteRequest.UserID))
	// if err!=nil {
	// 	log.Println("Error in fetching like status")
	// 	return
	// }

	if (prevVote == voteRequest.VoteType) {
		return;	
	}
	likeCountBuffer.mu.Lock()
	defer likeCountBuffer.mu.Unlock()
	if(voteRequest.VoteType == Like) {
		likeCountBuffer.hmap[voteRequest.PostID] = LikeCountBufferValues{likeCnt: likeCountBuffer.hmap[voteRequest.PostID].likeCnt + 1, dislikeCnt: likeCountBuffer.hmap[voteRequest.PostID].dislikeCnt + prevVote}
	}
	if(voteRequest.VoteType == Dislike) {
		likeCountBuffer.hmap[voteRequest.PostID] = LikeCountBufferValues{likeCnt: likeCountBuffer.hmap[voteRequest.PostID].likeCnt - prevVote, dislikeCnt: likeCountBuffer.hmap[voteRequest.PostID].dislikeCnt + 1}
	}
	if(voteRequest.VoteType == NotVoted) {
		if(prevVote == Like) {
			likeCountBuffer.hmap[voteRequest.PostID] = LikeCountBufferValues{likeCnt: likeCountBuffer.hmap[voteRequest.PostID].likeCnt - 1, dislikeCnt: likeCountBuffer.hmap[voteRequest.PostID].dislikeCnt}
		}
		if(prevVote == Dislike) {
			likeCountBuffer.hmap[voteRequest.PostID] = LikeCountBufferValues{likeCnt: likeCountBuffer.hmap[voteRequest.PostID].likeCnt, dislikeCnt: likeCountBuffer.hmap[voteRequest.PostID].dislikeCnt - 1}
		}
	}
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

	batchUpdaterDB, err := openDBConn()
	if err != nil {
		log.Print("Error in connecting to likes DB:\n", err)
		panic(err)
	}

	defer ch.BatchUpdater.Close()
	processingInterval := time.Duration(viper.GetInt("MQ.PROCESSING_INTERVAL")) * time.Second
	timer := time.NewTicker(processingInterval)

	var likeCountBuffer LikeCountBuffer
	likeCountBuffer.hmap = make(map[uuid.UUID]LikeCountBufferValues)
	var buffer []amqp.Delivery

	msgs, err := openBatchUpdaterQueue()
	if err != nil {
		return
	}

	defer processBufferedMessages(buffer, &likeCountBuffer, batchUpdaterDB)
	defer acknowledgeMessages(buffer)
	for {
		select {
		case message := <-msgs:
			// likeCountBuffer.updatelikeCountBuffer(message,batchUpdaterDB)
			buffer = append(buffer, message)
		case <-timer.C:
			processBufferedMessages(buffer, &likeCountBuffer, batchUpdaterDB)
			acknowledgeMessages(buffer)
			buffer = nil
		}
	}
}

func processBufferedMessages(buffer []amqp.Delivery, likeCountBuffer *LikeCountBuffer, batchUpdaterDB *gorm.DB) {
	if len(buffer) == 0 {
		return
	}
	logBuffer(buffer)

	likeCountBuffer.mu.Lock()
	defer likeCountBuffer.mu.Unlock()
	for pid, newLikeCnt := range likeCountBuffer.hmap {
		err := updateBatchLikeCount(batchUpdaterDB, pid, newLikeCnt)
		if err != nil {
			log.Printf("Error in updating likes count for pid: %d \n", pid)
		}
	}
	likeCountBuffer.hmap = make(map[uuid.UUID]LikeCountBufferValues)
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
