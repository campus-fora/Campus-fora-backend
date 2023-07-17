package users

import (
	"github.com/campus-fora/constants"
	"gorm.io/gorm"
)

type UserDetails struct {
	gorm.Model
	UserID     uint
	Department string `json:"department"`
	Program    string `json:"program"`
	Year       string `json:"year"`
	Gender     string `json:"gender"`
	Hall       string `json:"hall"`
	Hometown   string `json:"hometown"`
	UserBio    string `json:"user_bio"`
}

type User struct {
	gorm.Model
	UserId             string               `json:"user_id"`
	Email              string               `json:"email"`
	Password           string               `json:"password"`
	Role               constants.Role       `json:"role_id"`
	Name               string               `json:"name"`
	Blocked            bool                 `gorm:"default:false" json:"blocked" `
	LastLogin          uint                 `gorm:"index;autoUpdateTime:milli" json:"last_login"`
	ProfilePic         string               `json:"profile_pic"`
	UserDetails        UserDetails          `gorm:"foreignKey:UserID"`
	UserQuestions      []UserQuestions      `gorm:"foreignKey:UserID"` // ID is the foreign key
	UserLikedQuestions []UserLikedQuestions `gorm:"foreignKey:UserID"`
	Notifications      []Notification       `gorm:"foreignKey:UserID"`
	NotifTokens        []NotifTokens        `gorm:"foreignKey:UserID"`
}

type TemporaryUser struct {
	ID               uint   `gorm:"primary_key;autoIncrement:true" json:"user_id"`
	Name             string `gorm:"type:varchar(255);not null" json:"name"`
	Email            string `gorm:";not null" json:"email"`
	Password         string `gorm:"not null" json:"password"`
	Verificationcode string `gorm:"uniqueIndex;not null" json:"verificationcode"`
	Expires          uint   `gorm:"column:expires"`
}

type UserQuestions struct {
	gorm.Model
	UserID     uint
	QuestionId uint
}

type UserLikedQuestions struct {
	gorm.Model
	UserID     uint
	QuestionId uint
}

type Notification struct {
	gorm.Model
	UserID uint
	Title  string `json:"title"`
	Body   string `json:"body"`
	Read   bool   `json:"read"`
	Link   string `json:"link"`
}

type NotifTokens struct {
	gorm.Model
	UserID   uint
	Token    string `json:"token"`
	DeviceId string `json:"device_id"`
}
