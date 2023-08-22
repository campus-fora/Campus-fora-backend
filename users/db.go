package users

import "github.com/gin-gonic/gin"

func GetUserNameByID(ctx *gin.Context, userId uint) (error, string) {
	return nil, "user"
}
