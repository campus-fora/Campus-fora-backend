package likes

import (
	"time"
)

const (
	Like     = 1
	Dislike  = -1
	NotVoted = 0
)

type UserLike struct {
	UserID        uint `gorm:"index:user_like_idx;primaryKey" json:"user_id"`
	PostID        uint `gorm:"index:user_like_idx;primaryKey" json:"post_id"`
	VoteType      int  `json:"vote_type"`
	LatestReqTime time.Time
}

type TotalLikeCount struct {
	PostID       uint `gorm:"index:total_like_idx;primaryKey"`
	LikeCount    int
	DislikeCount int
}

type DailyLikeCount struct {
	PostID       uint `gorm:"post_id;primaryKey"`
	LikeCount    int
	DislikeCount int
}

type newVoteRequest struct {
	PostID        uint `json:"post_id"`
	UserID        uint `json:"user_id"`
	VoteType      int  `json:"vote_type"`
	LatestReqTime time.Time
}
