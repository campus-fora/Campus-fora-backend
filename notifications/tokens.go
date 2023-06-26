package notifications

import (
	"log"
	"net/http"

	"github.com/campus-fora/users"
	"github.com/gin-gonic/gin"
)

type newNotificationToken struct {
	notifToken string
	userId     uint
	deviceId   string
}

type deleteTokenRequest struct {
	userId   uint
	deviceId string
}

func postToken(ctx *gin.Context) {
	var notificationTokenRequest newNotificationToken

	if err := ctx.ShouldBindJSON(&notificationTokenRequest); err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := users.SaveNewToken(ctx, notificationTokenRequest.userId, notificationTokenRequest.notifToken, notificationTokenRequest.deviceId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification token saved successfully"})
}

func deleteToken(ctx *gin.Context) {
	var deleteTokenRequest deleteTokenRequest

	if err := ctx.ShouldBindJSON(&deleteTokenRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := users.DeleteToken(ctx, deleteTokenRequest.userId, deleteTokenRequest.deviceId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification token deleted successfully"})
}
