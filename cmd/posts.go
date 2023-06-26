package main

import (
	"net/http"
	"time"
	"github.com/spf13/viper"
	"github.com/campus-fora/middleware"
	"github.com/campus-fora/posts"
	"github.com/gin-gonic/gin"
)

func postsServer() *http.Server {
	READTIMEOUT := time.Duration(viper.GetInt("SERVER.READTIMEOUT")) * time.Second
	WRITETIMEOUT := time.Duration(viper.GetInt("SERVER.WRITETIMEOUT")) * time.Second
	PORT := viper.GetString("PORT.POSTS")
	engine := gin.Default()
	engine.Use(middleware.CORS())

	posts.InitRouters(engine)

	server := &http.Server{
		Addr:         ":"+PORT,
		Handler:      engine,
		ReadTimeout:  READTIMEOUT,
		WriteTimeout: WRITETIMEOUT,
	}

	return server
}
