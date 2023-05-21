package posts

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

func initCache() {
	redis_client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
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
	host := "database"
	port := "5432"
	password := "postkaadmin"
	dbName := "posts"
	user := "postsadmin"

	dsn := "host=" + host + " user=" + user + " password=" + password
	dsn += " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Print("Error in opening connection to posts DB:\n", err)
		panic(err)
	}

	db = database

	err = db.AutoMigrate(&Thread{}, &Tag{}, &Post{}, &Comment{})
	if err != nil {
		fmt.Print("Error in automigrating posts DB:\n", err)
		panic(err)
	}

	fmt.Println("Successfully connected to posts DB")
}

func Init() {
	initDB()
	initCache()
}
