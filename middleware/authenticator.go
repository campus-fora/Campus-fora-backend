package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"campus-fora/initializers"
	"time"

	"github.com/google/uuid"

)
func GetUserId(ctx *gin.Context) uint {   //dummy function
	return 0
}
func Authenticator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		config, _ := initializers.LoadConfig(".")
		user_id, err := ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
        type User struct {
			ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"user_id"`
			Name      string    `gorm:"type:varchar(255);not null" json:"name"`
			Email     string    `gorm:"uniqueIndex;not null" json:"email"`
			Password  string    `gorm:"not null" json:"password"`
			Role      string    `gorm:"type:varchar(255);not null" json:"role_id" `
			
			CreatedAt time.Time
			UpdatedAt time.Time
		}
		var user User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(user_id))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("user_id", user_id)
		ctx.Next()
	}
}
func GetUserID(ctx *gin.Context) string {
	return ctx.GetString("user_id")
}

