package posts

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func fetchAllUsers()

// func fetchUserDetailsByID()

// func updateUserDetails()

// func fetchUserQuestions()

// func fetchUserStarredQuestions()

func fetchUserStarQuestionStatus(ctx *gin.Context, userId uint, qid uint) (bool, error) {

	tx := db.WithContext(ctx).Model(&UserStarredQuestions{}).Where("user_id = ? and question_id = ?", userId, qid).First(&UserStarredQuestions{})
	starred := false
	if tx.RowsAffected > 0 {
		starred = true
	}
	return starred, tx.Error
}

func toggleOrCreateStarQuestion(ctx *gin.Context, userId uint, qid uint) error {
	var starredQuestion *UserStarredQuestions
	tx := db.WithContext(ctx).Model(&UserStarredQuestions{}).Unscoped().Where("user_id = ? and question_id = ?", userId, qid).First(starredQuestion)
	if(errors.Is(tx.Error, gorm.ErrRecordNotFound)){
		tx = db.WithContext(ctx).Model(&UserStarredQuestions{}).Create(&UserStarredQuestions{UserID: userId, QuestionId: qid})
		return tx.Error
	}
	if (starredQuestion.DeletedAt.Valid) {
		tx = tx.Update("deleted_at", gorm.Expr("NULL"))
	} else {
		tx = tx.Delete(starredQuestion)
	}
	return tx.Error
}

func fetchAllStarredQuestionsByUser(ctx *gin.Context, userId uint) ([]UserStarredQuestions, error) {
	var starredQuestions []UserStarredQuestions
	tx := db.WithContext(ctx).Model(&Question{}).Preload("UserStarredQuestions").Preload("Tags").Where("user_id = ?", userId).Find(&starredQuestions)
	return starredQuestions, tx.Error
}

// func fetchUserLikedQuestions()

// func fetchUserNotifications()

// func addUserQuestion()

// func deleteUserQuestion()

// func addUserStarOnQuestion()

// func deleteUserStarOnQuestion()

// func addUserLikeOnQuestion()

// func deleteUserLikeOnQuestion()

// func firstOrCreateNotification()

// func updateUserNotificationReadStatus()

// func deleteUserNotificationById()
