package users

import (
	"log"
	_ "github.com/campus-fora/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	host := "database"
	port := "5432"
	password := "userkaadmin"
	dbName := "users"
	user := "usersadmin"

	dsn := "host=" + host + " user=" + user + " password=" + password
	dsn += " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error in opening connection to users DB:\n", err)
		panic(err)
	}

	db = database

	err = db.AutoMigrate(&UserDetails{}, &User{}, &UserQuestions{}, &NotifTokens{}, &Notification{})
	if err != nil {
		log.Fatal("Error in automigrating users DB:\n", err)
		panic(err)
	}

	log.Println("Successfully connected to users DB")
}

func init() {
	initDB()
}
