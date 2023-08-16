package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func createComment(ctx *gin.Context, comment *Comment) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Create(comment)
	return tx.Error
}

func updateCommentByUUID(ctx *gin.Context, cid uuid.UUID, comment *Comment) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Where("uuid = ?", cid).Updates(comment)
	return tx.Error
}

func deleteCommentByUUID(ctx *gin.Context, cid uuid.UUID) error {
	tx := db.WithContext(ctx).Model(&Comment{}).Where("uuid = ?", cid).Delete(&Comment{})
	return tx.Error
}
