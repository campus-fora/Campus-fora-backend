package auth

import (
	"github.com/campus-fora/users"
	"github.com/campus-fora/constants"
)

type User = users.User

type TemporaryUser = users.TemporaryUser

type SignUpRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
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
}
