package stats

import (
	"gorm.io/gorm"
)

type User1 struct {
	gorm.Model
	UserID string `gorm:"uniqueIndex;not null" json:"user_id"`
	//Questions []Ques `gorm:"many2many:user_roles;"`
}

type Ques1 struct {
	gorm.Model
	QuesID string `gorm:"uniqueIndex;not null" json:"ques_id"`
	//Users  []User    `gorm:"many2many:user_roles;" json:"users"`
}
type User struct {
	gorm.Model
	UserID    string `gorm:"uniqueIndex;not null" json:"user_id"`
	Questions []Ques `gorm:"many2many:user_ques;" json:"questions"`
}

type Ques struct {
	gorm.Model
	QuesID string `gorm:"uniqueIndex;not null" json:"ques_id"`
	Users  []User `gorm:"many2many:user_ques;" json:"users"`
}

type UserQues struct {
	UserID string `gorm:"primaryKey"`
	QuesID string `gorm:"primaryKey"`
}
