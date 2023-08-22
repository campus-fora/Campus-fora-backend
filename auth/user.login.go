package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/campus-fora/middleware"
	"github.com/gin-gonic/gin"

	//"github.com/campus-fora/config"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *SignInInput

	AccessTokenExpiresIn, _ := time.ParseDuration(viper.GetString("JWT.EXPIRATION.ACCESS_TOKEN"))
	RefreshTokenExpiresIn, _ := time.ParseDuration(viper.GetString("JWT.EXPIRATION.REFRESH_TOKEN"))
	AccessTokenPrivateKey := viper.GetString("JWT.TOKENS.ACCESS_TOKEN_PRIVATE_KEY")
	RefreshTokenPrivateKey := viper.GetString("JWT.TOKENS.REFRESH_TOKEN_PRIVATE_KEY")
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	access_token, err := middleware.CreateToken(AccessTokenExpiresIn, user.ID, uint(user.Role), AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := middleware.CreateToken(RefreshTokenExpiresIn, user.ID, uint(user.Role), RefreshTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, int(AccessTokenExpiresIn), "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, int(RefreshTokenExpiresIn), "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", int(AccessTokenExpiresIn), "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}
