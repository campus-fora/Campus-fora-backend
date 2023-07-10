package likes

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(router *gin.Engine) {
	likesRouter := router.Group("api/likes")
	{
		likesRouter.PUT("/:pid", updateUserLikeStatusHandler)
		likesRouter.GET("/:pid", getUserLikeStatusHandler)
		likesRouter.GET("/:pid/count", getLikesCountHandler)
		likesRouter.GET("/:pid/likedQuestions", getLikedQuestionsByUser)
	}
}
