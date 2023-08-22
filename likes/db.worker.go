package likes

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func updateUserLikeStatus(db *gorm.DB, postId uuid.UUID, userId uint, voteType int, latestReqTime time.Time) (error , int, bool) {
	userLike := &UserLike{PostID: postId, UserID: userId}
	tx := db.Model(&UserLike{}).Where("post_id = ? AND user_id = ?", postId, userId).FirstOrCreate(userLike)
	if(userLike.LatestReqTime.After(latestReqTime)){
		return nil, userLike.VoteType, false
	}
	tx = tx.Updates(map[string]interface{}{"vote_type": voteType,"latest_req_time": latestReqTime})

	return tx.Error, userLike.VoteType, true
}