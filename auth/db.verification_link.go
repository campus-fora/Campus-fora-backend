package auth

import (
	"time"
	_ "time"

	_ "github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
)

func (ac *AuthController) cleanupOTP() {
	for {
		ac.DB.Unscoped().Delete(TemporaryUser{}, "expires < ?", time.Now().Add(-24*time.Hour).UnixMilli())
		time.Sleep(time.Hour * 24)
	}
}
