package posts

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type allQuestionResponse struct {
	UUID           uuid.UUID      `json:"uuid"`
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
		log.Print("cache hit")
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
	log.Print(err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func createNewQuestionHandler(ctx *gin.Context) {
	var question *Question

	if err := ctx.ShouldBindJSON(question); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		print(err.Error())
		return
	}
	question = &Question{UUID: uuid.New()}
	if err := createQuestion(ctx, question); err != nil {
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

func getQuestionHandler(ctx *gin.Context) {
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err!=nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}
	question := &Question{UUID: qid}
	err = fetchQuestionByUUID(ctx, question)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, question)
}

func updateQuestionHandler(ctx *gin.Context) {
	var question Question;
	if err := ctx.ShouldBindJSON(&question); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	qid := ctx.Param("qid")

	if err := updateQuestion(ctx, qid, &question); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, question)

}

func updateQuestionTagsHandler(ctx *gin.Context) {
	var tags []Tag;
	if err := ctx.ShouldBindJSON(&tags); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	qid := ctx.Param("qid")

	if err := updateQuestionTags(ctx, qid, &tags); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, tags)
}

// func getQuestionsByTagsHandler(ctx *gin.Context) {
// 	qid := ctx.Param("qid")
// 	tags := []Tag{}
// 	err := fetchQuestionTags(ctx, qid, &tags)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return;
// 	}
// 	ctx.JSON(http.StatusOK, tags)
// }
func deleteQuestionHandler(ctx *gin.Context) {
	qid := ctx.Param("qid")
	if err := deleteQuestionByUUID(ctx, qid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

func createAnswerHandler(ctx *gin.Context) {
	var answer Answer;
	if err := ctx.ShouldBindJSON(&answer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}
	answer.UUID = uuid.New()
	qid := answer.ParentID
	
	if err := createAnswer(ctx,qid ,&answer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	
	ctx.JSON(http.StatusOK, answer)
}
func updateAnswerHandler(ctx *gin.Context) {
	var answer Answer;
	if err := ctx.ShouldBindJSON(&answer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	aid := ctx.Param("aid")
	if err := updateAnswerByUUID(ctx, aid, &answer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, answer)
}

func deleteAnswerHandler(ctx *gin.Context) {
	aid := ctx.Param("aid")
	if err := deleteAnswerByUUID(ctx, aid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully"})
}
func createCommentHandler(ctx *gin.Context) {
	var comment Comment;
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}
	comment.UUID = uuid.New()

	if err := createComment(ctx, &comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}

	ctx.JSON(http.StatusOK, comment)
}

func updateCommentHandler(ctx *gin.Context) {
	var comment Comment;
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	cid := ctx.Param("cid")
	if err := updateCommentByUUID(ctx, cid, &comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, comment)
}

func deleteCommentHandler(ctx *gin.Context) {
	cid := ctx.Param("cid")
	if err := deleteCommentByUUID(ctx, cid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}






