package posts

import (
	// "net/http"

	"github.com/campus-fora/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func updateQuestionFollowingStatus(ctx *gin.Context) {
	qid, err := strconv.ParseUint(ctx.Param("qid"), 10 ,32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := middleware.GetUserId(ctx)
	err = toggleOrCreateStarQuestion(ctx, userId, uint(qid))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "successfully updated question following status"})
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
	qid, err := strconv.ParseUint(ctx.Param("qid"), 10 ,32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	starred, err := fetchUserStarQuestionStatus(ctx, userId, uint(qid))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": starred})
}
