package posts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	Name      string     `json:"name"`
	Questions []Question `gorm:"foreignKey:TopicID"`
	Tags      []Tag      `gorm:"foreignKey:TopicID"`
}

type Question struct {
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt         `gorm:"index"`
	UUID                 uuid.UUID              `json:"uuid" gorm:"primaryKey"`
	TopicID              uint                   `json:"topic_id"`
	Title                string                 `json:"title"`
	Content              string                 `json:"content"`
	CreatedByUser        uint                   `json:"created_by_user"`
	CreatedByUserName    string                 `json:"created_by_user_name"`
	LikeCount            LikeCount              `gorm:"foreignKey:PostID"`
	Answers              []Answer               `gorm:"foreignKey:ParentID" json:"answers"`
	Tags                 []Tag                  `gorm:"many2many:question_tags;" json:"tags"`
	UserStarredQuestions []UserStarredQuestions `gorm:"foreignKey:QuestionId"`
	ReportedByUser       []ReportedByUser       `gorm:"foreignKey:QuestionId"`
}

type ReportedByUser struct {
	QuestionId uint `gorm:"primaryKey"`
	UserID     uint `gorm:"primaryKey"`
	Reason     string
}

type QuestionRelevancy struct {
	QuestionID uuid.UUID `json: "uuid" gorm:"primaryKey"`
	Relevancy  int       `json: "relevancy"`
}
type Tag struct {
	gorm.Model
	TopicID uint   `json:"topic_id"`
	Name    string `json:"name"`
}
type Answer struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	UUID              uuid.UUID      `json:"uuid" gorm:"primaryKey"`
	ParentID          uuid.UUID      `json:"parent_id"`
	Content           string         `json:"content"`
	IsAnswer          bool           `json:"is_answer"`
	CreatedByUser     uint           `json:"created_by_user"`
	CreatedByUserName string         `json:"created_by_user_name"`
	LikeCount         LikeCount      `gorm:"foreignKey:PostID"`
	Comments          []Comment      `gorm:"foreignKey:ParentID" json:"comments"`
}

type Comment struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	UUID              uuid.UUID      `json:"uuid" gorm:"primaryKey"`
	ParentID          uuid.UUID      `json:"parent_id"`
	Content           string         `json:"content"`
	CreatedByUser     uint           `json:"created_by_user"`
	CreatedByUserName string         `json:"created_by_user_name"`
}

type QuestionDetail struct {
	UUID          uuid.UUID `json:"uuid"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedByUser uint      `json:"created_by_user"`
}

type UserStarredQuestions struct {
	UserID     uint           `gorm:"primaryKey"`
	QuestionId uuid.UUID      `gorm:"primaryKey"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// type TotalLikeCount struct {
// 	PostID       uuid.UUID `gorm:"index:total_like_idx;primaryKey"`
// 	LikeCount    int
// 	DislikeCount int
// }

type LikeCount struct {
	PostID       uuid.UUID `gorm:"post_id;primaryKey"`
	LikeCount    int       `gorm:"like_count" json:"likeCount"`
	DislikeCount int       `gorm:"dislike_count" json:"dislikeCount"`
}
