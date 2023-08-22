package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"strconv"
)

func createTopicHandler(ctx *gin.Context){
	var topic Topic
	if err := ctx.ShouldBindJSON(&topic); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := createTopic(ctx, &topic); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, topic)
}

func getAllTopicsHandler(ctx *gin.Context){
	var topics []Topic
	err := fetchAllTopics(ctx, &topics)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, topics)
}

func updateTopicHandler(ctx *gin.Context){
	var topic Topic
	if err := ctx.ShouldBindJSON(&topic); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tid, err := strconv.ParseUint(ctx.Param("tid"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := updateTopic(ctx, uint(tid), &topic); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, topic)
}

func deleteTopicHandler(ctx *gin.Context){
	tid, err := strconv.ParseUint(ctx.Param("tid"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := deleteTopic(ctx, uint(tid)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "successfully deleted topic"})
}

func getAllTopicTagsHandler(ctx *gin.Context){
	tid, err := strconv.ParseUint(ctx.Param("tid"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tags []Tag
	err = fetchAllTopicTags(ctx, uint(tid), &tags)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tags)
}