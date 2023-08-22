package config

import (
	"log"

	"github.com/spf13/viper"
)

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
		panic(err)
	}
	log.Print("Successfully loaded config file")
}
