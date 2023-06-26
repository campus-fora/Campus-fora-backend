package posts

import (
	"time"

	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model
	Title         string `json:"title"`
	Content       string `json:"content"`
	CreatedByUser string `json:"created_by_user"`
	Posts         []Post `gorm:"foreignKey:Parent_ID" json:"posts"`
	Tags          []Tag  `gorm:"many2many:thread_tags;" json:"tags"`
}

type Tag struct {
	gorm.Model
	Name string `json:"name"`
}

type Post struct {
	gorm.Model
	Parent_ID     uint      `json:"parent_id"`
	Content       string    `json:"content"`
	IsAnswer      bool      `json:"is_answer"`
	CreatedByUser string    `json:"created_by_user"`
	Comments      []Comment `gorm:"foreignKey:Parent_ID"`
}

type Comment struct {
	gorm.Model
	Parent_ID     uint   `json:"parent_id"`
	Content       string `json:"content"`
	CreatedByUser string `json:"created_by_user"`
}

type thread_detail struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedByUser string    `json:"created_by_user"`
}
