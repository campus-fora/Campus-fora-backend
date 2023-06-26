package main

import (
	"net/http"
	"time"

	"github.com/campus-fora/middleware"
	"github.com/campus-fora/notifications"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func notificationServer() *http.Server {
	READTIMEOUT := time.Duration(viper.GetInt("SERVER.READTIMEOUT")) * time.Second
	WRITETIMEOUT := time.Duration(viper.GetInt("SERVER.WRITETIMEOUT")) * time.Second
	PORT := viper.GetString("PORT.NOTIFICATION")
	engine := gin.Default()
	engine.Use(middleware.CORS())
	notifications.InitRouters(engine)

	server := &http.Server{
		Addr:         ":" + PORT,
		Handler:      engine,
		ReadTimeout:  READTIMEOUT,
		WriteTimeout: WRITETIMEOUT,
	}

	return server
}
