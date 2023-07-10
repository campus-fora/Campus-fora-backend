package posts

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "gorm.io/gorm"
)

func fetchAllTopics(ctx *gin.Context, topics *[]Topic) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Find(topics)
	return tx.Error
}

func createTopic(ctx *gin.Context, topic *Topic) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Create(topic)
	return tx.Error
}

func updateTopic(ctx *gin.Context, tid uint, topic *Topic) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Where("id = ?", tid).Updates(Topic{Name: topic.Name})
	return tx.Error
}

func deleteTopic(ctx *gin.Context, tid uint) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Where("id = ?", tid).Delete(&Topic{})
	return tx.Error
}

func fetchAllQuestionDetails(ctx *gin.Context, tid uint, questionDetail *[]questionDetailResponse) error {
	// tx := db.WithContext(ctx).Model(&Question{}).Select("uuid, created_at, title, content, created_by_user").Preload("Tags").Find(questionDetail)
	// return tx.Error
	tx := db.WithContext(ctx).Model(&Question{}).Preload("Tags").Where("questions.topic_id = ?", tid).Joins("JOIN answers ON questions.uuid = answers.parent_id").Where("answers.deleted_at IS NULL").Select("questions.uuid, questions.created_at,questions.updated_at, questions.title, questions.content, questions.created_by_user, questions.created_by_user_name, string_agg(distinct answers.uuid::text, ',') as answer_ids").Group("questions.uuid").Find(questionDetail)
	return tx.Error
}

func fetchQuestionWithoutAnswer(ctx *gin.Context, qid uuid.UUID, question *Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Preload("Tags").Where("questions.uuid = ?", qid).First(question)
	return tx.Error
}

func fetchQuestionByUUID(ctx *gin.Context, question *Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Preload("Answers").Preload("Tags").First(question)
	return tx.Error
}

func createQuestion(ctx *gin.Context, question *Question) error {
	return db.WithContext(ctx).Model(&Question{}).FirstOrCreate(question).Error
}

func fetchLimitedQuestionByRelevancy(ctx *gin.Context, offset int, pageSize int, questions *[]Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Joins("JOIN question_relevancy ON question.uuid = question_relevancy.uuid").Order("relevancy DESC").Offset(offset).Limit(pageSize).Find(&questions)
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
func updateQuestion(ctx *gin.Context, qid uuid.UUID, question *Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("uuid = ?", qid).Updates(Question{Title: question.Title, Content: question.Content})
	return tx.Error
}

func updateQuestionTags(ctx *gin.Context, qid uuid.UUID, tags *[]Tag) error {
	return db.WithContext(ctx).Model(&Question{}).Where("uuid = ?", qid).Association("Tags").Replace(tags)
}

func deleteQuestionByUUID(ctx *gin.Context, qid uuid.UUID) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("uuid = ?", qid).Delete(&Question{})
	return tx.Error
}

func FetchAllQuestionsWithID(ctx *gin.Context, questionIds []uint, questions *[]Question) error {
	tx := db.WithContext(ctx).Model(&Question{}).Where("uuid IN ?", questionIds).Find(questions)
	return tx.Error
}

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

func createComment(ctx *gin.Context, comment *Comment) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Create(comment)
	return tx.Error
}

func updateCommentByUUID(ctx *gin.Context, cid uuid.UUID, comment *Comment) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Where("uuid = ?", cid).Updates(comment)
	return tx.Error
}

func deleteCommentByUUID(ctx *gin.Context, cid uuid.UUID) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Where("uuid = ?", cid).Delete(&Comment{})
	return tx.Error
}
func QuestionExists(ques_id uint64) bool {
	var ques Question
	result := db.First(&ques, "ID = ?", fmt.Sprint(ques_id))
	return result.Error == nil

}
