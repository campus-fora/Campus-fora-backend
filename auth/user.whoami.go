package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"campus-fora/middleware"
	"campus-fora/initializers"

	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}
func fetchUser(ctx *gin.Context, user *User, userID string) error {
	tx := initializers.DB.WithContext(ctx).Where("ID = ?", userID).First(&user)
	return tx.Error
}
func whoamiHandler(ctx *gin.Context) {
    middleware.Authenticator()(ctx)
	user_id := middleware.GetUserID(ctx)
	if user_id == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Null user id"})
		return
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

