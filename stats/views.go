package stats

import (
	"fmt"
	"log"
	"net/http"

	"github.com/campus-fora/auth"
	"github.com/campus-fora/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateViewHandler(ctx *gin.Context, qid uuid.UUID) error {

	middleware.Authenticator()(ctx)
	user_id := middleware.GetUserID(ctx)
	if user_id == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Null user id"})
		return nil
	}
	if !auth.UserExists(user_id) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Authorization Provided"})

		return nil
	}

	err := updateViewList(qid, user_id)
	return err
	//updateViewCount(tid, user_id)

}
func updateViewList(qid uuid.UUID, user_id string) error {
	var ques1 Ques1
	var user User1

	err := DB.Where("ques_id = ?", fmt.Sprint(qid.String())).First(&ques1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ques1.QuesID = qid.String()
			result := DB.Create(&ques1)
			if result.Error != nil {
				log.Print("Error updating stats Queslist")
				//return result.Error
			}
			log.Print("Stats Queslist updated")
		}
	}

	err = DB.Where("user_id = ?", user_id).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user.UserID = user_id
			result := DB.Create(&user)
			if result.Error != nil {
				log.Print("Error updating stats Userlist")
			}
			log.Print("Stats Userlist updated")
		}
	}
	var userQues UserQues

	err = DB.Where("user_id = ? AND ques_id = ?", user_id, qid.String()).First(&userQues).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userQues.UserID = user_id
			userQues.QuesID = qid.String()
			result := DB.Create(&userQues)
			if result.Error != nil {
				log.Print("Error updating stats Userlist Relation")
			}
			log.Print("Stats Userlist updated Relation")
		}
	}

	return err

}
