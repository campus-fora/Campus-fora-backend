package users

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(router *gin.Engine) {
	users := router.Group("/api/users")
	{
		users.GET("/", getUserDetailHandler)
		users.PUT("/", updateUserDetailHandler)
		users.GET("/questions", getUserQuestionsHandler)
		users.GET("/liked", getUserLikedQuestionsHandler)
		// 	// users.GET("/:id/notifications", getUserNotificationHandler)
		// 	// users.DELETE("/:id/notifications/:notificationId", deleteUserNotificationHandler)
		// 	// users.PUT("/:id/notifications/:notificationId", markUserNotificationAsReadHandler)
		// 	//
	}
}
