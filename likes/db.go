package likes

import (
	// "log"

	// "context"
	// "log"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func updateBatchLikeCount(db *gorm.DB, postId uuid.UUID, newVoteCount LikeCountBufferValues) error {

	tx := db.Model(&DailyLikeCount{}).Where("post_id = ?", postId).FirstOrCreate(&DailyLikeCount{PostID: postId}).
    Updates(map[string]interface{}{
        "like_count":    gorm.Expr("COALESCE(like_count, 0) + ?", newVoteCount.likeCnt),
        "dislike_count": gorm.Expr("COALESCE(dislike_count, 0) + ?", newVoteCount.dislikeCnt),
    })

	return tx.Error
}

func fetchLikesCountForPost(ctx *gin.Context, postId uuid.UUID) (DailyLikeCount, error) {
	var dailyCount DailyLikeCount
	// var totalCount TotalLikeCount
	tx := db.WithContext(ctx).Model(&DailyLikeCount{}).FirstOrCreate(&dailyCount, DailyLikeCount{PostID: postId})
	if(tx.Error != nil){
		return DailyLikeCount{}, tx.Error
	}
	// tx = db.WithContext(ctx).Model(&TotalLikeCount{}).FirstOrCreate(&totalCount, TotalLikeCount{PostID: postId})
	// if(tx.Error != nil){
	// 	return DailyLikeCount{}, tx.Error
	// }
	return DailyLikeCount{PostID: postId, LikeCount: dailyCount.LikeCount, DislikeCount: dailyCount.DislikeCount}, tx.Error
}

func fetchLikeStatus(ctx *gin.Context, postId uuid.UUID, userId uint) (int, error) {
	var userLike UserLike
	tx := db.WithContext(ctx).Model(&UserLike{}).FirstOrCreate(&userLike,&UserLike{UserID: userId, PostID: postId})
	return userLike.VoteType, tx.Error
}
func fetchLikeStatusWithoutContext(db *gorm.DB,postId uint, userId uint) (int, error ){
	userLike := &UserLike{}
	tx := db.WithContext(context.Background()).Model(&UserLike{}).Where("post_id = ? AND user_id = ?", postId, userId).First(userLike)
	return userLike.VoteType, tx.Error
}
func fetchLikedPostsByUser(ctx *gin.Context, userId uint, allPostIds *[]uint) error {
	tx := db.WithContext(ctx).
		Model(&UserLike{}).
		Where("user_id = ?", userId).
		Select("post_id").
		Find(allPostIds)
	return tx.Error
}

