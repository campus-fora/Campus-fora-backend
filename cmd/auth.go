package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "github.com/campus-fora/config"
	"github.com/campus-fora/mail"
	"github.com/campus-fora/middleware"

	"github.com/campus-fora/auth"
)

var (
	server              *gin.Engine
	AuthController      auth.AuthController
	AuthRouteController auth.AuthRouteController

	UserController      auth.UserController
	UserRouteController auth.UserRouteController
)

func init() {

	AuthController = auth.NewAuthController(auth.DB)
	AuthRouteController = auth.NewAuthRouteController(AuthController)

	UserController = auth.NewUserController(auth.DB)
	UserRouteController = auth.NewRouteUserController(UserController)

	server = gin.Default()
}

func authServer(mail_channel chan mail.Mail) *http.Server {
	READTIMEOUT := time.Duration(viper.GetInt("SERVER.READTIMEOUT")) * time.Second
	WRITETIMEOUT := time.Duration(viper.GetInt("SERVER.WRITETIMEOUT")) * time.Second
	PORT := viper.GetString("PORT.LIKES")
	engine := gin.Default()
	engine.Use(middleware.CORS())
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	router := engine.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Auth Service Working..."
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	go mail.Service(mail_channel)

	AuthRouteController.AuthRoute(mail_channel, router)
	UserRouteController.UserRoute(router)
	log.Fatal(server.Run(":" + "8000"))
	server := &http.Server{
		Addr:         ":" + PORT,
		Handler:      engine,
		ReadTimeout:  READTIMEOUT,
		WriteTimeout: WRITETIMEOUT,
	}

	return server
}
