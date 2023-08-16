package posts

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/campus-fora/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func fetchAllQuestionDetails(ctx *gin.Context, tid uint, questionDetail *[]QuestionDetailResponse) error {
	tx := db.WithContext(ctx).Model(&Question{}).Preload("Tags").Where("questions.topic_id = ?", tid).Joins("JOIN answers ON questions.uuid = answers.parent_id").Where("answers.deleted_at IS NULL").Select("questions.uuid, questions.created_at,questions.updated_at, questions.title, questions.content, questions.created_by_user, questions.created_by_user_name, string_agg(distinct answers.uuid::text, ',') as answer_ids").Group("questions.uuid").Find(questionDetail)
	return tx.Error
}

func fetchFilteredAndPaginatedQuestions(ctx *gin.Context, topicID uint, sortBy string, tags []uint, lastKeySetValue string, limit int, questions *[]QuestionDetailResponse) error {
	tx := db.Model(&Question{}).
		Where("questions.topic_id = ?", topicID).
		Joins("JOIN question_tags ON questions.uuid = question_tags.question_uuid").
		Joins("JOIN tags ON tags.id = question_tags.tag_id").
		Joins("LEFT JOIN like_counts ON questions.uuid = like_counts.post_id").
		Joins("LEFT JOIN answers ON questions.uuid = answers.parent_id").
		Where("tags.id IN (?)", tags).
		Select("questions.uuid, questions.created_at,questions.updated_at, questions.title, questions.content, questions.created_by_user, questions.created_by_user_name, string_agg(distinct answers.uuid::text, ',') as answer_ids, SUM(like_counts.like_count) as like_count, SUM(like_counts.dislike_count) as dislike_count,json_agg(json_build_object('tag_id', tags.id, 'tag_name', tags.name)) as tags").
		Group("questions.uuid").
		Having("COUNT(DISTINCT tags.id) = ?", len(tags))

	switch sortBy {
	case "newest":
		if lastKeySetValue != "" {
			parsedKeySet, err := utils.ParseTime(lastKeySetValue)
			if err != nil {
				return err
			}
			tx = tx.Order("created_at DESC").Where("questions.created_at < ?", parsedKeySet).Limit(limit).Find(&questions)
		} else {
			tx = tx.Order("created_at DESC").Limit(limit).Find(&questions)
		}
	case "upvotes":
		if lastKeySetValue != "" {
			parsedKeySet, err := strconv.ParseUint(lastKeySetValue, 10, 32)
			if err != nil {
				return err
			}
			tx = tx.Order("like_count DESC").Where("like_counts.like_count < ?", parsedKeySet).Limit(limit).Find(&questions)
		} else {
			tx = tx.Order("like_count DESC").Limit(limit).Find(&questions)
		}
	default:
		return errors.New("Invalid sort by value")
	}
	return tx.Error
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

func QuestionExists(ques_id uint64) bool {
	var ques Question
	result := db.First(&ques, "ID = ?", fmt.Sprint(ques_id))
	return result.Error == nil
}
