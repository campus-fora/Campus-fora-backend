package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/campus-fora/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (ac *AuthController) VerificationLinkHandler(mail_channel chan mail.Mail) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Params.ByName("verificationCode")
		var verifyuser TemporaryUser
		tx := ac.DB.WithContext(ctx).Where("verificationcode = ? AND expires > ?", code, time.Now().UnixMilli()).First(&verifyuser)
		switch tx.Error {
		case nil:
			ac.DB.WithContext(ctx).Delete(&verifyuser)
			//now := time.Now()
			newUser := User{
				Name:     verifyuser.Name,
				Email:    strings.ToLower(verifyuser.Email),
				Password: verifyuser.Password,
				Role:     1,

				//CreatedAt: now,
				//UpdatedAt: now,
			}
			result := ac.DB.Create(&newUser)

			if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
				ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
				return
			} else if result.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened", "error": result.Error})
				return
			}
			mail_channel <- mail.GenerateMail(newUser.Email, "Registered on Campus-Fora", "Dear "+newUser.Name+",\n\nYou have been registered on Campus-Fora")
			ctx.JSON(http.StatusOK, gin.H{"status": "Successfully signed up"})

		case gorm.ErrRecordNotFound:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Verification Link may have expired, Try again."})
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Something bad happened", "error": tx.Error})

			return
		}

	}
}
