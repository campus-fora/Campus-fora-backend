package posts

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(router *gin.Engine) {
	posts := router.Group("/")
	{
		posts.GET("threads", getAllThreadsDetail)
		posts.POST("postThread", postThread)
	}
}
