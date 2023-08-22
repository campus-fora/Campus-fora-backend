package users

import (
	// "errors"
	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func fetchUserDetails(ctx *gin.Context, userId uint) (UserDetails, error) {
	var userDetails UserDetails
	err := Db.WithContext(ctx).Model(&UserDetails{}).Where("user_id=?", userId).First(&userDetails)
	return userDetails, err.Error
}

func updateUserDetails(ctx *gin.Context, userId uint, userDetails UserDetails) error {
	return Db.WithContext(ctx).Model(&UserDetails{}).Where("user_id=?", userId).Updates(&userDetails).Error
}

func fetchUserQuestions(ctx *gin.Context, userId uint) ([]UserQuestions, error) {
	var userQuestions []UserQuestions
	err := Db.WithContext(ctx).Model(&UserQuestions{}).Where("user_id=?", userId).Find(&userQuestions)
	return userQuestions, err.Error
}

func fetchUserLikedQuestions(ctx *gin.Context, userId uint) ([]UserLikedQuestions, error) {
	var userLikedQuestions []UserLikedQuestions
	err := Db.WithContext(ctx).Model(&UserLikedQuestions{}).Where("user_id=?", userId).Find(&userLikedQuestions)
	return userLikedQuestions, err.Error
}
