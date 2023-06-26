package likes

import (
	"gorm.io/gorm"
)

const (
	Like     = 1
	Dislike  = -1
	NotVoted = 0
)

type UserLike struct {
	gorm.Model
	UserID   uint `json:"user_id"`
	PostID   uint `json:"post_id"`
	VoteType int  `json:"vote_type"`
}

type TotalLikeCount struct {
	gorm.Model
	PostID uint `json:"post_id"`
	Count  uint `json:"count"`
}

type DailyLikeCount struct {
	gorm.Model
	PostID uint `json:"post_id"`
	Count  uint `json:"count"`
}

type newVoteRequest struct {
	PostID   uint `json:"post_id"`
	UserID   uint `json:"user_id"`
	VoteType int  `json:"vote_type"`
}
