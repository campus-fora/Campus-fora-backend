package posts

import (
	"net/http"

	"github.com/campus-fora/middleware"
	"github.com/campus-fora/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NewCommentReq struct {
	ParentID          uuid.UUID      `json:"parentID"`
	Content           string         `json:"content"`
}

func createCommentHandler(ctx *gin.Context) {
	var commentReq NewCommentReq
	if err := ctx.ShouldBindJSON(&commentReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := middleware.GetUserId(ctx)
	err, userName := users.GetUserNameByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	comment := Comment{
		UUID:          uuid.New(),
		ParentID:      commentReq.ParentID,
		Content:       commentReq.Content,
		CreatedByUser: userId,
		CreatedByUserName: userName,
	}

	if err := createComment(ctx, &comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func updateCommentHandler(ctx *gin.Context) {
	var comment Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cid, err := uuid.Parse(ctx.Param("cid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	comment = Comment{
		Content: comment.Content,
	}
	if err := updateCommentByUUID(ctx, cid, &comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comment)
}

func deleteCommentHandler(ctx *gin.Context) {
	cid , err := uuid.Parse(ctx.Param("cid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if err := deleteCommentByUUID(ctx, cid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
