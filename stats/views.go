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

func UpdateViewHandler(ctx *gin.Context, qid uuid.UUID) {
	middleware.Authenticator()(ctx)
	user_id := middleware.GetUserID(ctx)
	fmt.Print("yo_user")
	if user_id == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Null user id"})
		return
	}
	if !auth.UserExists(user_id) {
		return
	}
	updateViewList(qid, user_id)
	//updateViewCount(tid, user_id)

}
func updateViewList(qid uuid.UUID, user_id string) {
	var ques Ques

	// DB.FirstOrCreate(&ques, Ques{QuesID: qid})
	// if DB.NewRecord(ques) {
	// 	ques.Users = []User{
	// 		{UserID: user_id},
	// 			}
	// 	DB.Create(&ques)
	// } else {

	// 	ques.Users = append(ques.Users, User{UserID: user_id})
	// 	DB.Save(&ques)
	// }
	// if DB.Model(&ques).Where("ques_id ?", qid).Updates(&user).RowsAffected == 0 {
	// 	db.Create(&user)
	// }
	if err := DB.Model(&ques).Where("ques_id = ?", qid).Update("user", append(ques.Users, User{UserID: user_id})).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ques.QuesID = qid
			ques.Users = []User{
				{UserID: user_id},
			}
			result := DB.Create(&ques)
			if result.Error != nil {
				log.Print("Error updating viewlist")
				return
			}
			fmt.Print("Viewlist updated")
		}
	}

}
