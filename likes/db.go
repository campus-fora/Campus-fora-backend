package likes

import (
	// "log"

	// "context"
	// "log"

	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func updateBatchLikeCount(db *gorm.DB, postId uint, newVoteCount LikeCountBufferValues) error {

	tx := db.Model(&DailyLikeCount{}).Where("post_id = ?", postId).FirstOrCreate(&DailyLikeCount{PostID: postId}).
    Updates(map[string]interface{}{
        "like_count":    gorm.Expr("COALESCE(like_count, 0) + ?", newVoteCount.likeCnt),
        "dislike_count": gorm.Expr("COALESCE(dislike_count, 0) + ?", newVoteCount.dislikeCnt),
    })

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
func fetchLikeStatusWithoutContext(postId uint, userId uint) (int, error ){
	// log.Print("inside fetchLikeStatusWithoutContext", postId, userId)
	// var userLike *UserLike
	userLike := &UserLike{}
	tx := db.WithContext(context.Background()).Model(&UserLike{}).Where("post_id = ? AND user_id = ?", postId, userId).First(userLike)
	// log.Println(" debugging ---> ")
	return userLike.VoteType, tx.Error
	// return 0, nil
}
func fetchLikedPostsByUser(ctx *gin.Context, userId uint, allPostIds *[]uint) error {
	tx := db.WithContext(ctx).
		Model(&UserLike{}).
		Where("user_id = ?", userId).
		Select("post_id").
		Find(allPostIds)
	return tx.Error
}

