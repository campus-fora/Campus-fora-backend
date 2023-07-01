package likes

import (
	"log"

	_ "github.com/campus-fora/config"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQChannels struct {
	Publisher    *amqp.Channel
	Updater      *amqp.Channel
	BatchUpdater *amqp.Channel
}

var ch MQChannels
var db *gorm.DB
var voteCh chan newVoteRequest

func openDBConn() (*gorm.DB, error) {
	host := viper.GetString("DATABASE.HOST")
	port := viper.GetString("DATABASE.PORT")
	password := "likekaadmin"
	dbName := viper.GetString("DBNAME.LIKES")
	user :=  viper.GetString("DBNAME.LIKES") + viper.GetString("DATABASE.USER")

	dsn := "host=" + host + " user=" + user + " password=" + password
	dsn += " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return database, nil
}

func initDB() {
	db, err := openDBConn()
	if err != nil {
		log.Print("Error in connecting to likes DB:\n", err)
		panic(err)
	}

	err = db.AutoMigrate(&UserLike{}, &TotalLikeCount{}, &DailyLikeCount{})
	if err != nil {
		log.Print("Error in automigrating likes DB:\n", err)
		panic(err)
	}

	log.Println("Successfully connected to likes DB")
	log.Println(db)
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
}

func init() {
	initDB()
	initMQ()
	voteCh = make(chan newVoteRequest)

	go publisher()
	go updater()
	// go batchUpdater()
}