package posts

import (
	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func fetchAllQuestionDetails(ctx *gin.Context, questionDetail *[]Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Select("id, created_at, title, content, created_by_user").Preload("Tags").Find(questionDetail)
	return tx.Error
}

func createQuestion(ctx *gin.Context, question *Question) error {
	tx := db.WithContext(ctx).Create(question)
	return tx.Error
}

func FetchAllQuestionsWithID(ctx *gin.Context, questionIds []uint, questions *[]Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("id IN ?", questionIds).Find(questions)
	return tx.Error
}

func FetchAllAnswersWithId(ctx *gin.Context, answerIds []uint, answers *[]Question) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("id IN ?", answerIds).Find(&[]Answer{})
	return tx.Error
}