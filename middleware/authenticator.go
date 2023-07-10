package middleware

import (
	//"fmt"
	"net/http"
	"strings"

	"github.com/campus-fora/constants"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	//"time"
	//"github.com/google/uuid"
)

func GetUserId(ctx *gin.Context) uint { //dummy function
	return 0
}
func Authenticator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		AccessTokenPublicKey := viper.GetString("JWT.TOKENS.ACCESS_TOKEN_PUBLIC_KEY")

		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "authorization header is not provided"})
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "invalid authorization header format"})
			return
		}
		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		user_id, role_id, err := ValidateToken(access_token, AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message2": err.Error()})
			return
		}

		ctx.Set("user_id", user_id)
		ctx.Set("role_id", role_id)

		ctx.Next()
	}
}
// func GetUserID(ctx *gin.Context) string {
// 	return ctx.GetString("user_id")
// }
func GetRoleID(ctx *gin.Context) constants.Role {
	return constants.Role(ctx.GetInt("role_id"))
}
