package auth

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	host := "127.0.0.1"
	port := viper.GetString("DATABASE.PORT")
	password := "postgres"
	dbName := "postgres"
	user := "postgres"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", host, user, password, dbName, port)
	fmt.Print("\nyohoho" + dsn + "\n")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = database
	if err != nil {
		fmt.Print("Failed to connect to the User DB\n", err)
		panic(err)
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		fmt.Print("Error in automigrating user DB:\n", err)
		panic(err)
	}
	err = DB.AutoMigrate(&TemporaryUser{})
	if err != nil {
		fmt.Print("Error in automigrating Temporary User DB:\n", err)
		panic(err)
	}

	fmt.Println("Successfully connected to user DB")
}

func init() {
	initDB()
}
