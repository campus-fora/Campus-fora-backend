package auth

import (
	"time"

	"github.com/campus-fora/constants"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint           `gorm:"primary_key;autoIncrement:true" json:"user_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"password"`
	Role      constants.Role `gorm:"default:1;not null" json:"role_id" ` //IITK_USER by default
	Blocked   bool           `gorm:"default:false" json:"blocked" `
	LastLogin uint           `gorm:"index;autoUpdateTime:milli" json:"last_login"`

	//CreatedAt time.Time
	//UpdatedAt time.Time
}

type SignUpRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID    uint           `json:"id,omitempty"`
	Name  string         `json:"name,omitempty"`
	Email string         `json:"email,omitempty"`
	Role  constants.Role `json:"role,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type TemporaryUser struct {
	ID               uint   `gorm:"primary_key;autoIncrement:true" json:"user_id"`
	Name             string `gorm:"type:varchar(255);not null" json:"name"`
	Email            string `gorm:";not null" json:"email"`
	Password         string `gorm:"not null" json:"password"`
	Verificationcode string `gorm:"uniqueIndex;not null" json:"verificationcode"`
	Expires          uint   `gorm:"column:expires"`
}
