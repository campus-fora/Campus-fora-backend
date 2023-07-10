package main

import (
	"log"

	_ "github.com/campus-fora/config"
	// "github.com/campus-fora/mail"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.Print("Starting server")
	var g errgroup.Group
	// mail_channel := make(chan mail.Mail)

	gin.SetMode(gin.ReleaseMode)

	g.Go(func() error {
		return postsServer().ListenAndServe()
	})
	// g.Go(func() error {
	// 	return authServer(mail_channel).ListenAndServe()
	// })

	g.Go(func() error {
		return notificationServer().ListenAndServe()
	})

	g.Go(func() error {
		return likesServer().ListenAndServe()
	})

	err := g.Wait()
	if err != nil {
		panic(err)
	}
}
