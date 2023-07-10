package posts

import (
	"log"
	"net/http"

	"github.com/campus-fora/middleware"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func updateQuestionFollowingStatus(ctx *gin.Context) {
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("error parsing uuid", err)
		return
	}

	userId := middleware.GetUserId(ctx)
	err, status := toggleOrCreateStarQuestion(ctx, userId, qid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, status)
}

func getAllStarredQuestions(ctx *gin.Context) {
	// userId := middleware.GetUserId(ctx)
	// starredQuestionsByUser, err := fetchAllStarredQuestionsByUser(ctx, userId)

	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

}

func getQuestionFollowingStatus(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	starred, err := fetchFollowingStatus(ctx, userId, qid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, starred)
}
