package likes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func updateUserLike(postId uint, userId uint, voteType int) error {
	tx := db.Model(&UserLike{}).FirstOrCreate(&UserLike{}, &UserLike{PostID: postId, UserID: userId, VoteType: voteType})
	return tx.Error
}

func updateBatchVoteCount(postId uint, netVoteCount int) error {
	tx := db.Model(&DailyLikeCount{}).Where("post_id = ?", postId).Update("count", gorm.Expr("count + ?", netVoteCount))
	return tx.Error
}

func fetchLikesCountForPost(ctx *gin.Context, postId uint) (uint, error) {
	var totalLikes uint
	tx := db.WithContext(ctx).Model(&DailyLikeCount{}).Joins("daily_like_counts.post_id = total_like_counts.post_id").Select("daily_like_counts.count + total_like_counts.count").Find(&totalLikes)
	return totalLikes, tx.Error
}

func fetchLikeStatus(ctx *gin.Context, postId uint, userId uint) (int, error) {
	var userLike UserLike
	tx := db.WithContext(ctx).Model(&UserLike{}).Where("post_id = ? AND user_id = ?", postId, userId).First(&userLike)
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
