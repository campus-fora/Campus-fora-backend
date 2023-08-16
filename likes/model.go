package likes

import (
	"time"

	"github.com/campus-fora/posts"
	"github.com/google/uuid"
)

const (
	Like     = 1
	Dislike  = -1
	NotVoted = 0
)

type DailyLikeCount = posts.LikeCount

type UserLike struct {
	UserID        uint `gorm:"index:user_like_idx;primaryKey" json:"user_id"`
	PostID        uuid.UUID `gorm:"index:user_like_idx;primaryKey" json:"post_id"`
	VoteType      int  `json:"vote_type"`
	LatestReqTime time.Time
}

type newVoteRequest struct {
	PostID        uuid.UUID `json:"post_id"`
	UserID        uint `json:"user_id"`
	VoteType      int  `json:"vote_type"`
	LatestReqTime time.Time
}
