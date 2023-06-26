package likes

import (
	"fmt"
	"log"
	
	_ "github.com/campus-fora/config"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQChannels struct {
	Publisher    *amqp.Channel
	Updater      *amqp.Channel
	BatchUpdater *amqp.Channel
}

var ch RabbitMQChannels
var db *gorm.DB
var voteCh chan newVoteRequest

func initDB() {
	host := viper.GetString("DATABASE.HOST")
	port := viper.GetString("DATABASE.PORT")
	password := viper.GetString("DATABASE.PASSWORD")
	dbName := viper.GetString("DATABASE.DBNAME")
	user := viper.GetString("DATABASE.USER") + viper.GetString("DATABASE.DBNAME")

	dsn := "host=" + host + " user=" + user + " password=" + password
	dsn += " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Print("Error in opening connection to posts DB:\n", err)
		panic(err)
	}

	db = database

	err = db.AutoMigrate(&UserLike{}, &TotalLikeCount{}, &DailyLikeCount{})
	if err != nil {
		log.Print("Error in automigrating posts DB:\n", err)
		panic(err)
	}

	log.Println("Successfully connected to posts DB")
}

func openMQChannels(conn *amqp.Connection) error {
	var err error
	ch.Publisher, err = conn.Channel()
	if err != nil {
		log.Printf("Error in establishing publisher Channel : %s", err)
		return err
	}

	ch.Updater, err = conn.Channel()
	if err != nil {
		log.Printf("Error in establishing updater Channel : %s", err)
		return err
	}

	ch.BatchUpdater, err = conn.Channel()
	if err != nil {
		log.Printf("Error in establishing batch updater Channel : %s", err)
		return err
	}
	return nil
}

func initMQ() {
	amqpURI := viper.GetString("MQ.URI")
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("Error in establishing AMQP Connection : %s", err)
		return
	}
	err = openMQChannels(conn)
	if err != nil {
		log.Printf("Error in opening AMQP Channels : %s", err)
		return
	}

	go publisher()
	go updater()
	go batchUpdater()
}

func processBatch(batch []amqp.Delivery) {
	fmt.Println("Processing batch:")
	for _, msg := range batch {
		fmt.Printf("Message: %s\n", msg.Body)
	}
	fmt.Println()
}

func init() {
	// initDB()
	initMQ()
	voteCh = make(chan newVoteRequest)
}
