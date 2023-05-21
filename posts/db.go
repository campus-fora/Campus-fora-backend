package posts

import (
	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func fetchAllThreadDetails(ctx *gin.Context, thread_detail *[]Thread) error {
	tx := db.WithContext(ctx).Model(&Thread{}).Select("id, created_at, title, content, created_by_user").Preload("Tags").Find(thread_detail)
	return tx.Error
}

func createThread(ctx *gin.Context, thread *Thread) error {
	tx := db.WithContext(ctx).Create(thread)
	return tx.Error
}
