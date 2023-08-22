package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/campus-fora/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")
	AccessTokenPrivateKey := viper.GetString("JWT.TOKENS.ACCESS_TOKEN_PRIVATE_KEY")

	//RefreshTokenPrivateKey:= ""
	RefreshTokenPublicKey := viper.GetString("JWT.TOKENS.REFRESH_TOKEN_PUBLIC_KEY")

	AccessTokenExpiresIn, _ := time.ParseDuration(viper.GetString("JWT.EXPIRATION.ACCESS_TOKEN"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	userID, _, err := middleware.ValidateToken(cookie, RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(userID))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := middleware.CreateToken(AccessTokenExpiresIn, user.ID, uint(user.Role), AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, int(AccessTokenExpiresIn), "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", int(AccessTokenExpiresIn), "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}
