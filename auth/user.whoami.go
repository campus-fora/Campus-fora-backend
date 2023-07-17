package auth

import (
	"fmt"
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
	fmt.Print("yo_user")
	if user_id == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Null user id2"})
		return
	}
	if !UserExists(user_id) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token does not exist"})

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
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}
