package main

import (
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

	//UserController auth.UserController
	//UserRouteController auth.UserRouteController
)

func init() {

	AuthController = auth.NewAuthController(auth.DB)
	AuthRouteController = auth.NewAuthRouteController(AuthController)

	//UserController = auth.NewUserController(auth.DB)
	//UserRouteController = auth.NewRouteUserController(UserController)

}

func authServer(mail_channel chan mail.Mail) *http.Server {
	READTIMEOUT := time.Duration(viper.GetInt("SERVER.READTIMEOUT")) * time.Second
	WRITETIMEOUT := time.Duration(viper.GetInt("SERVER.WRITETIMEOUT")) * time.Second
	PORT := viper.GetString("PORT.AUTH")
	engine := gin.Default()
	engine.Use(middleware.CORS())

	go mail.Service(mail_channel)

	AuthRouteController.AuthRoute(mail_channel, engine)
	//UserRouteController.UserRoute(engine)
	server := &http.Server{
		Addr:         ":" + PORT,
		Handler:      engine,
		ReadTimeout:  READTIMEOUT,
		WriteTimeout: WRITETIMEOUT,
	}

	return server
}
