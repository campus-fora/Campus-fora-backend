package likes

import "gorm.io/gorm"

func updateUserLikeStatus(db *gorm.DB, postId uint, userId uint, voteType int) error {
	tx := db.Model(&UserLike{}).Where("post_id = ? AND user_id = ?", postId, userId).FirstOrCreate(&UserLike{PostID: postId, UserID: userId}).Updates(map[string]interface{}{"vote_type": voteType,})
	return tx.Error
}