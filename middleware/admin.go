package middleware

import (
	"net/http"

	"github.com/campus-fora/constants"
	"github.com/gin-gonic/gin"
)

func EnsureAuthority() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := GetRoleID(ctx)

		if role != constants.ADMIN && role != constants.MOD {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
