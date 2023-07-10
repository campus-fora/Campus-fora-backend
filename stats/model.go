package stats

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID    string `gorm:"uniqueIndex;not null"`
	Questions []Ques `gorm:"many2many:user_roles;"`
}

type Ques struct {
	gorm.Model
	QuesID uuid.UUID `gorm:"uniqueIndex;not null" json:"ques_id"`
	Users  []User    `gorm:"many2many:user_roles;" json:"user"`
}
