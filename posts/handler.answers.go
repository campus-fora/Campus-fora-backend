package posts

import (
	"net/http"

	"github.com/campus-fora/middleware"
	"github.com/campus-fora/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getAnswerByUUIDHandler(ctx *gin.Context) {
	var answer Answer
	aid, err := uuid.Parse(ctx.Param("aid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := fetchAnswerByUUID(ctx, aid, &answer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, answer)
}

func getAllAnswersWithUUIDsHandler(ctx *gin.Context) {
	var answers []Answer
	var answerIds []uuid.UUID
	// aid, err := uuid.Parse(ctx.Param("aid"))
	for _, aid := range ctx.QueryArray("aid") {
		aid, err := uuid.Parse(aid)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		answerIds = append(answerIds, aid)
	}
	if err := fetchAnswersWithUUIDs(ctx, answerIds, &answers); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, answers)
}

func createAnswerHandler(ctx *gin.Context) {
	var answer *Answer
	if err := ctx.ShouldBindJSON(&answer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := middleware.GetUserId(ctx)
	err, userName := users.GetUserNameByID(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	answer = &Answer{
		UUID:          uuid.New(),
		ParentID:      answer.ParentID,
		Content:       answer.Content,
		CreatedByUser: userId,
		CreatedByUserName: userName,
	}

	if err := createAnswer(ctx, answer, userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, answer)
}

func updateAnswerHandler(ctx *gin.Context) {
	var answer Answer
	if err := ctx.ShouldBindJSON(&answer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aid, err := uuid.Parse(ctx.Param("aid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	answer = Answer{
		Content: answer.Content,
	}

	if err := updateAnswerByUUID(ctx, aid, &answer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, answer)
}

func deleteAnswerHandler(ctx *gin.Context) {
	aid, err := uuid.Parse(ctx.Param("aid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := deleteAnswerByUUID(ctx, aid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully"})
}

func updateIsAnswerCorrectHandler(ctx *gin.Context) {
	type requestBody struct {
		IsAnswer bool `json:"is_answer"`
	}
	var body requestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	aid, err := uuid.Parse(ctx.Param("aid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	answer := Answer{
		IsAnswer: body.IsAnswer,
	}
	if err := updateAnswerByUUID(ctx, aid, &answer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "successfully changed answer's status"})
}
