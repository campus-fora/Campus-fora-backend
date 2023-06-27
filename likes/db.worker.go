package likes

import "gorm.io/gorm"

func updateUserLikeStatus(db *gorm.DB, postId uint, userId uint, voteType int) error {
	tx := db.Model(&UserLike{}).FirstOrCreate(&UserLike{}, &UserLike{PostID: postId, UserID: userId, VoteType: voteType})
	return tx.Error
}