package auth

import (
	//"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/campus-fora/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (ac *AuthController) signUpHandler(mail_channel chan mail.Mail) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload *SignUpRequest

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := HashPassword(payload.Password)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
			return
		}

		//now := time.Now()
		// newUser := User{
		// 	Name:     payload.Name,
		// 	Email:    strings.ToLower(payload.Email),
		// 	Password: hashedPassword,
		// 	Role:     1,

		// 	CreatedAt: now,
		// 	UpdatedAt: now,
		// }
		var newUser User

		tx := ac.DB.WithContext(ctx).Where("email = ?", strings.ToLower(payload.Email)).First(&newUser)
		if tx.Error != gorm.ErrRecordNotFound {
			if tx.Error == nil {
				ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists, Try a new mail or Contact: "})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened", "error": tx.Error})
			return
		}
		var secret_code = generateCode()
		tempUser := TemporaryUser{
			Name:             payload.Name,
			Email:            strings.ToLower(payload.Email),
			Password:         hashedPassword,
			Verificationcode: secret_code,
			Expires:          uint(time.Now().Add(time.Duration(linkExpiration) * time.Minute).UnixMilli()),
		}

		result := ac.DB.Create(&tempUser)

		if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {

			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "fail", "message": "Not able to resolve unique verification code conflict"})
			return
		} else if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened", "error": result.Error})
			return
		}

		// userResponse := &UserResponse{
		// 	ID:    newUser.ID,
		// 	Name:  newUser.Name,
		// 	Email: newUser.Email,

		// 	CreatedAt: newUser.CreatedAt,
		// 	UpdatedAt: newUser.UpdatedAt,
		// }
		mail_channel <- mail.GenerateMail(tempUser.Email, "eefeefeff", "Dear "+newUser.Name+",\n\n"+"code "+secret_code+"\n")
		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Check your mail inbox to verify"})
	}
}
