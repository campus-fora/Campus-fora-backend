package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func fetchAnswerByUUID(ctx *gin.Context, aid uuid.UUID, answer *Answer) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Preload("Comments").First(answer, aid)
	return tx.Error
}

func fetchAnswersWithUUIDs(ctx *gin.Context, aid []uuid.UUID, answers *[]Answer) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Preload("Comments").Where("UUID in ?", aid).Find(answers)
	return tx.Error
}

func FetchAllAnswersWithUUID(ctx *gin.Context, answerIds []uint, answers *[]Question) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("uuid IN ?", answerIds).Find(&[]Answer{})
	return tx.Error
}

func createAnswer(ctx *gin.Context, answer *Answer, userId uint) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Create(answer)
	return tx.Error
}

func updateAnswerByUUID(ctx *gin.Context, aid uuid.UUID, answer *Answer) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("uuid = ?", aid).Updates(answer)
	return tx.Error
}

func deleteAnswerByUUID(ctx *gin.Context, aid uuid.UUID) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("uuid = ?", aid).Delete(&Answer{})
	return tx.Error
}
