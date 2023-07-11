package stats

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	host := viper.GetString("DATABASE.HOST")
	port := viper.GetString("DATABASE.PORT")
	password := "statskaadmin"
	dbName := viper.GetString("DBNAME.STATS")
	user := "statsadmin"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", host, user, password, dbName, port)
	fmt.Print("\nyohoho" + dsn + "\n")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = database
	if err != nil {
		fmt.Print("Failed to connect to the Analytics DB\n", err)
		panic(err)
	}

	err = DB.AutoMigrate(&User1{}, &Ques1{}, &UserQues{})
	if err != nil {
		fmt.Print("Error in automigrating Analytics DB:\n", err)
		panic(err)
	}

	fmt.Println("Successfully connected to Analytics DB")
}

func init() {
	initDB()
}
