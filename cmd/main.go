package main

import (
	"campus-fora/posts"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	posts.InitRouters(r)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})
	posts.Init()

	err := g.Wait()
	if err != nil {
		panic(err)
	}
}
