package posts

import (
	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func fetchAllQuestionDetails(ctx *gin.Context, questionDetail *[]Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Select("uuid, created_at, title, content, created_by_user").Preload("Tags").Find(questionDetail)
	return tx.Error
}

func fetchQuestionByUUID(ctx *gin.Context, question *Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Preload("Answers").Preload("Tags").First(question)
	return tx.Error
}

func createQuestion(ctx *gin.Context, question *Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Create(question)
	return tx.Error
}

func fetchQuestionsByTag(ctx *gin.Context, tags []string, questions *[]Question) error {
	tx := db.Preload("Tags").
    Joins("JOIN question_tags ON questions.uuid = question_tags.question_uuid").
    Joins("JOIN tags ON tags.id = question_tags.tag_id").
    Where("tags.name IN (?)", tags).
    Group("questions.uuid").
    Having("COUNT(DISTINCT tags.id) = ?", len(tags)).
    Find(&questions)
	return tx.Error
}
func updateQuestion(ctx *gin.Context, qid string ,question *Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("uuid = ?", qid).Updates(question)
	return tx.Error
}

func updateQuestionTags(ctx *gin.Context, qid string, tags *[]Tag) error {
	return db.WithContext(ctx).Model(&Question{}).Where("uuid = ?", qid).Association("Tags").Replace(tags)
}

func deleteQuestionByUUID(ctx *gin.Context, qid string) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("uuid = ?", qid).Delete(&Question{})
	return tx.Error
}

func FetchAllQuestionsWithID(ctx *gin.Context, questionIds []uint, questions *[]Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("uuid IN ?", questionIds).Find(questions)
	return tx.Error
}

func FetchAllAnswersWithId(ctx *gin.Context, answerIds []uint, answers *[]Question) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("uuid IN ?", answerIds).Find(&[]Answer{})
	return tx.Error
}

func createAnswer(ctx *gin.Context, qid string, answer *Answer) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Create(answer)
	return tx.Error
}

func updateAnswerByUUID(ctx *gin.Context, aid string, answer *Answer) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("uuid = ?", aid).Updates(answer)
	return tx.Error
}

func deleteAnswerByUUID(ctx *gin.Context, aid string) error {
	tx := db.WithContext(ctx).Model(&Answer{}).Where("uuid = ?", aid).Delete(&Answer{})
	return tx.Error
}

func createComment(ctx *gin.Context, comment *Comment) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Create(comment)
	return tx.Error
}

func updateCommentByUUID(ctx *gin.Context, cid string, comment *Comment) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Where("uuid = ?", cid).Updates(comment)
	return tx.Error
}

func deleteCommentByUUID(ctx *gin.Context, cid string) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Where("uuid = ?", cid).Delete(&Comment{})
	return tx.Error
}



