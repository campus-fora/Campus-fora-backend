package posts

import (
	"context"
	"fmt"

	_ "github.com/campus-fora/config"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

func initCache() {
	redis_client := redis.NewClient(&redis.Options{
		Addr:     "redis:" + viper.GetString("REDIS.PORT"),
		Password: "",
		DB:       0,
	})

	if _, err := redis_client.Ping(context.Background()).Result(); err != nil {
		fmt.Print("Error connecting to Redis:", err.Error())
		rdb = nil
		return
	}

	rdb = redis_client
	fmt.Println("Successfully connected to Redis")
}

func initDB() {
	host := viper.GetString("DATABASE.HOST")
	port := viper.GetString("DATABASE.PORT")
	password := "postkaadmin"
	dbName := viper.GetString("DBNAME.POSTS")
	user := "postsadmin"

	dsn := "host=" + host + " user=" + user + " password=" + password
	dsn += " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Print("Error in opening connection to posts DB:\n", err)
		panic(err)
	}

	db = database

	err = db.AutoMigrate(&Topic{} , &Question{}, &Tag{}, &Answer{}, &Comment{}, &UserStarredQuestions{})
	if err != nil {
		fmt.Print("Error in automigrating posts DB:\n", err)
		panic(err)
	}

	fmt.Println("Successfully connected to posts DB")
}

func init() {
	initDB()
	initCache()
}
