package posts

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// func fetchAllUsers()

// func fetchUserDetailsByID()

// func updateUserDetails()

// func fetchUserQuestions()

// func fetchUserStarredQuestions()

func fetchFollowingStatus(ctx *gin.Context, userId uint, qid uuid.UUID) (bool, error) {

	tx := db.WithContext(ctx).Model(&UserStarredQuestions{}).FirstOrCreate(&UserStarredQuestions{UserID: userId, QuestionId: qid})
	starred := false
	if tx.RowsAffected > 0 {
		starred = true
	}
	if(errors.Is(tx.Error, gorm.ErrRecordNotFound)){
		return starred, nil
	}
	if tx.Error != nil {
		return starred, tx.Error
	}
	return starred, nil
}

func toggleOrCreateStarQuestion(ctx *gin.Context, userId uint, qid uuid.UUID) (error, bool) {
	var starredQuestion UserStarredQuestions
	var status = false
	tx := db.WithContext(ctx).Model(&UserStarredQuestions{}).Unscoped().Where("user_id = ? and question_id = ?", userId, qid).First(&starredQuestion)
	if(errors.Is(tx.Error, gorm.ErrRecordNotFound)){
		log.Print("record not found")
		tx = db.WithContext(ctx).Model(&UserStarredQuestions{}).Create(&UserStarredQuestions{UserID: userId, QuestionId: qid})
		status = true
		return tx.Error, status
	}
	log.Print("record found")
	if (starredQuestion.DeletedAt.Valid) {
		log.Print("record found and deleted")
		tx = tx.Update("deleted_at", gorm.Expr("NULL"))
		status = true
	} else {
		log.Print("record found and not deleted")
		tx = tx.Delete(starredQuestion)
	}
	return tx.Error, status
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
