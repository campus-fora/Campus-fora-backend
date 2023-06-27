package likes

const (
	Like     = 1
	Dislike  = -1
	NotVoted = 0
)

type UserLike struct {
	UserID   uint `gorm:"column:user_id;index:user_like_idx;primaryKey"`
	PostID   uint `gorm:"column:post_id;index:user_like_idx;primaryKey"`
	VoteType int  `gorm:"vote_type"`
}

type TotalLikeCount struct {
	PostID       uint `gorm:"column:post_id";"index:total_like_idx";"primaryKey"`
	LikeCount    int  `gorm:"like_count"`
	DislikeCount int  `gorm:dislike_count`
}

type DailyLikeCount struct {
	PostID       uint `gorm:"post_id;primaryKey"`
	LikeCount    int  `gorm:"like_count"`
	DislikeCount int  `gorm:dislike_count`
}

type newVoteRequest struct {
	PostID   uint `json:"post_id"`
	UserID   uint `json:"user_id"`
	VoteType int  `json:"vote_type"`
}
