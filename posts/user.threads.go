package posts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type allQuestionResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedByUser string    `json:"created_by_user"`
	Tags          []Tag     `json:"tags"`
}

func getAllQuestionsDetailHandler(ctx *gin.Context) {
	var questions []allQuestionResponse

	err := getAllQuestionDetailsCache(ctx, &questions)
	if err == nil {
		fmt.Print("cache hit")
		ctx.JSON(http.StatusOK, questions)
		return
	}
	if err == redis.Nil {
		var questions []Question
		err = fetchAllQuestionDetails(ctx, &questions)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		go setAllQuestionDetailCache(ctx, questions)
		ctx.JSON(http.StatusOK, questions)
		return
	}
	fmt.Print(err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func createNewQuestionHandler(ctx *gin.Context) {
	var question Question

	if err := ctx.ShouldBindJSON(&question); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		print(err.Error())
		return
	}

	if err := createQuestion(ctx, &question); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go setQuestionCache(ctx, question)

	ctx.JSON(http.StatusOK, question)
}

func getPosts(ctx *gin.Context) {
	response := struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}{
		Question: "What is the meaning of life?",
		Answer:   "42",
	}

	ctx.JSON(http.StatusOK, response)
}

// func getQuestionHandler(ctx *gin.Context) {
// 	qid := 
// }

func createAnswerHandler(ctx *gin.Context) {
	var answer Answer;

	if err := ctx.ShouldBindJSON(&answer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
}