package auth

import (
	"net/http"

	"github.com/campus-fora/middleware"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func whoamiHandler(ctx *gin.Context) {
	middleware.Authenticator()(ctx)
	user_id := middleware.GetUserID(ctx)
	if user_id == "" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Null user id"})
		return
	}
	if !UserExists(user_id) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "user does not exist"})
	}

	var currentUser User
	err := fetchUser(ctx, &currentUser, user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	userResponse := &UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Role:      currentUser.Role,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}
