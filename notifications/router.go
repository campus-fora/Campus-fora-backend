package notifications

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(router *gin.Engine) {
	notifRouter := router.Group("notif")
	{
		notifRouter.POST("/postToken", postToken)
		notifRouter.POST("/deleteToken", deleteToken)
		notifRouter.POST("/test", TestNotification)
	}
}
