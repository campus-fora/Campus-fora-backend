package likes

import (
	"log"
	"net/http"
	"time"

	"github.com/campus-fora/middleware"
	"github.com/campus-fora/posts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Vote struct {
	Value int `json:"value"`
}

func updateUserLikeStatusHandler(ctx *gin.Context) {
	var vote Vote

	userId := middleware.GetUserId(ctx) // fetch userId from middleware

	pid, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&vote); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if vote.Value > 1 || vote.Value < -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote value"})
		return
	}

	voteCh <- newVoteRequest{
		PostID:        pid,
		UserID:        userId,
		VoteType:      vote.Value,
		LatestReqTime: time.Now(),
	}
	log.Println("Vote vlaue ->", vote.Value)
	ctx.JSON(http.StatusOK, gin.H{"message": "Vote request recieved"})
}

func getLikesCountHandler(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
	}

	totalCount, err := fetchLikesCountForPost(ctx, pid)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vote count"})
		return
	}

	ctx.JSON(http.StatusOK, totalCount)
}

func getLikedQuestionsByUser(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)

	var allQuestionIds []uint
	err := fetchLikedPostsByUser(ctx, userId, &allQuestionIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var LiekdQuestions []posts.Question

	err = posts.FetchAllQuestionsWithID(ctx, allQuestionIds, &LiekdQuestions)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, LiekdQuestions)
}

func getUserLikeStatusHandler(ctx *gin.Context) {
	// get user id from context
	var userId = 0 // hardcoding
	postId, err := uuid.Parse(ctx.Param("pid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}
	voteStatus, err := fetchLikeStatus(ctx, postId, uint(userId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, voteStatus)
}
