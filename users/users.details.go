package users

import (
	"log"
	"net/http"

	"github.com/campus-fora/middleware"
	"github.com/gin-gonic/gin"
)

func getUserDetailHandler(ctx *gin.Context) {
	// create a DB call in users DB and fetch the user details from the user model where the user id matches
	userId := middleware.GetUserId(ctx)
	userDetails, err := fetchUserDetails(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("error fetching user details", err)
		return
	}
	ctx.JSON(http.StatusOK, userDetails)
}

func updateUserDetailHandler(ctx *gin.Context) {
	// create a DB call in users DB and update the user details from the user model where the user id matches
	userId := middleware.GetUserId(ctx)
	var userDetails UserDetails
	err := ctx.BindJSON(&userDetails)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("error binding user details", err)
		return
	}
	err = updateUserDetails(ctx, userId, userDetails)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("error updating user details", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user details updated successfully"})

}

func getUserQuestionsHandler(ctx *gin.Context) {
	// create a DB call in posts DB and fetch all questions from the question model where the created_by_user is the user id
	userId := middleware.GetUserId(ctx)
	questions, err := fetchUserQuestions(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("error fetching user questions", err)
		return
	}
	ctx.JSON(http.StatusOK, questions)
}

func getUserLikedQuestionsHandler(ctx *gin.Context) {
	// create a DB call in posts DB and fetch all questions from the user_starred_questions model where the user_id is the user id
	userId := middleware.GetUserId(ctx)
	questions, err := fetchUserLikedQuestions(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("error fetching user starred questions", err)
		return
	}
	ctx.JSON(http.StatusOK, questions)
}
