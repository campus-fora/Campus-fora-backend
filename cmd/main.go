package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	_ "github.com/campus-fora/config"
)

func main() {
	log.Print("Starting server")
	var g errgroup.Group
	gin.SetMode(gin.ReleaseMode)

	g.Go(func() error {
		return postsServer().ListenAndServe()
	})

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