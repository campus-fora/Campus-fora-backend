package posts

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type questionAsHmap struct {
	questionDetail []byte `redis:"question_detail"`
	answers        []byte `redis:"answers"`
	tags           []byte `redis:"tags"`
}

func getAllQuestionDetailsCache(ctx *gin.Context, questions *[]allQuestionResponse) error {
	cursor := uint64(0)
	size, err := rdb.DBSize(ctx).Result()
	if err != nil {
		return err
	}
	if size == 0 {
		return redis.Nil
	}

	for {
		keys, nextCursor, err := rdb.Scan(ctx, cursor, "question:*", 10).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			var questionDetail QuestionDetail
			var tags []Tag

			cmd, err := rdb.HGet(ctx, key, "question_detail").Bytes()
			if err != nil {
				return err
			}

			b := bytes.NewReader(cmd)

			if err := gob.NewDecoder(b).Decode(&questionDetail); err != nil && err.Error() != "EOF" {
				return err
			}

			cmd, err = rdb.HGet(ctx, key, "tags").Bytes()
			if err != nil {
				return err
			}

			b = bytes.NewReader(cmd)

			if err := gob.NewDecoder(b).Decode(&tags); err != nil && err.Error() != "EOF" {
				return err
			}

			question_detail_with_tags := allQuestionResponse{
				ID:            questionDetail.ID,
				CreatedAt:     questionDetail.CreatedAt,
				Title:         questionDetail.Title,
				Content:       questionDetail.Content,
				CreatedByUser: questionDetail.CreatedByUser,
				Tags:          tags,
			}
			*questions = append(*questions, question_detail_with_tags)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	fmt.Println("cache hit")
	return nil
}

func setQuestionCache(ctx *gin.Context, question Question) error {
	var bQuestionDetail, bAnswers, bTags bytes.Buffer

	questionDetail := QuestionDetail{
		ID:            question.ID,
		CreatedAt:     question.CreatedAt,
		Title:         question.Title,
		Content:       question.Content,
		CreatedByUser: question.CreatedByUser,
	}

	err := gob.NewEncoder(&bQuestionDetail).Encode(questionDetail)
	if err != nil {
		fmt.Print("error in cahcing: ", err.Error())
		return err
	}

	err = gob.NewEncoder(&bAnswers).Encode(question.Answers)
	if err != nil {
		return err
	}

	err = gob.NewEncoder(&bQuestionDetail).Encode(question.Tags)
	if err != nil {
		fmt.Print("error in cahcing: ", err.Error())
		return err
	}

	err = rdb.HSet(ctx, "question:"+strconv.FormatUint(uint64(question.ID), 10),
		"question_detail", bQuestionDetail.Bytes(),
		"answers", bAnswers.Bytes(),
		"tags", bTags.Bytes()).Err()
	if err != nil {
		fmt.Print("error in cahcing: ", err.Error())
	}
	fmt.Print("Caching completed")
	return err
}

func setQuestionDetailsCache(ctx *gin.Context, question Question) error {
	var b_question_detail, b_Tags bytes.Buffer

	question_detail := QuestionDetail{
		ID:            question.ID,
		CreatedAt:     question.CreatedAt,
		Title:         question.Title,
		Content:       question.Content,
		CreatedByUser: question.CreatedByUser,
	}

	err := gob.NewEncoder(&b_question_detail).Encode(question_detail)
	if err != nil {
		return err
	}

	err = gob.NewEncoder(&b_question_detail).Encode(question.Tags)
	if err != nil {
		return err
	}

	return rdb.HSet(ctx, "question:"+strconv.FormatUint(uint64(question.ID), 10),
		"question_detail", b_question_detail.Bytes(),
		"tags", b_Tags.Bytes()).Err()
}

func setAllQuestionDetailCache(ctx *gin.Context, questions []Question) error {
	for _, question := range questions {
		err := setQuestionDetailsCache(ctx, question)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("caching complete")
	return nil
}
