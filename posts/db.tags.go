package posts

import "github.com/gin-gonic/gin"

func createTag(ctx* gin.Context, tag* Tag) error {
	tx := db.WithContext(ctx).Model(&Tag{}).Create(tag)
	return tx.Error
}

func fetchTagsWithIds(ctx* gin.Context, tags* []Tag, tagIds []uint) error {
	tx := db.WithContext(ctx).Model(&Tag{}).Find(tags,tagIds)
	return tx.Error
}