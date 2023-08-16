package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addTagHandler(ctx* gin.Context){
	var tag Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !topicExists(ctx, tag.TopicID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "topic does not exist"})
		return
	}
	if err := createTag(ctx, &tag); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tag)
}