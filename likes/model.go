package likes

import (
	"time"

	"github.com/google/uuid"
)

const (
	Like     = 1
	Dislike  = -1
	NotVoted = 0
)

type UserLike struct {
	UserID        uint `gorm:"index:user_like_idx;primaryKey" json:"user_id"`
	PostID        uuid.UUID `gorm:"index:user_like_idx;primaryKey" json:"post_id"`
	VoteType      int  `json:"vote_type"`
	LatestReqTime time.Time
}

type TotalLikeCount struct {
	PostID       uuid.UUID `gorm:"index:total_like_idx;primaryKey"`
	LikeCount    int
	DislikeCount int
}

type DailyLikeCount struct {
	PostID       uuid.UUID `gorm:"post_id;primaryKey"`
	LikeCount    int  `gorm:"like_count" json:"likeCount"`
	DislikeCount int  `gorm:"dislike_count" json:"dislikeCount"`
}

type newVoteRequest struct {
	PostID        uuid.UUID `json:"post_id"`
	UserID        uint `json:"user_id"`
	VoteType      int  `json:"vote_type"`
	LatestReqTime time.Time
}
