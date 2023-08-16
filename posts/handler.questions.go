package posts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/campus-fora/middleware"
	"github.com/campus-fora/stats"
	"github.com/campus-fora/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type allQuestionResponse struct {
	UUID          uuid.UUID `json:"uuid"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedByUser uint      `json:"created_by_user"`
	Tags          []Tag     `json:"tags"`
}

type QuestionDetailResponse struct {
	UUID              uuid.UUID `json:"uuid"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Title             string    `json:"title"`
	Content           string    `json:"content"`
	CreatedByUser     uint      `json:"created_by_user"`
	CreatedByUserName string    `json:"created_by_user_name"`
	AnswerIds         string    `json:"answer_ids"`
	Tags              []Tag     `json:"tags" gorm:"many2many:question_tags;"`
	LikeCount         uint      `json:"like_count"`
	DislikeCount      uint      `json:"dislike_count"`
}

func getAllQuestionsDetailHandler(ctx *gin.Context) {
	var questions []QuestionDetailResponse

	// err := getAllQuestionDetailsCache(ctx, &questions)
	// if err == nil {
	// 	log.Print("cache hit")
	// 	ctx.JSON(http.StatusOK, questions)
	// 	return
	// }
	// if err == redis.Nil {
	// 	var questions []Question
	// 	err = fetchAllQuestionDetails(ctx, &questions)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	// go setAllQuestionDetailCache(ctx, questions)
	// 	ctx.JSON(http.StatusOK, questions)
	// 	return
	// }
	// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	tid, err := strconv.ParseUint(ctx.Param("tid"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !QuestionExists(tid) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Question does not exists"})
	}

	err = fetchAllQuestionDetails(ctx, uint(tid), &questions)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, questions)
}

func getPaginatedQuestionsHandler(ctx *gin.Context) {
	tid, err := strconv.ParseUint(ctx.Param("tid"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sortBy := ctx.Query("sortBy")
	tags := ctx.QueryArray("tag")
	var tagIds []uint
	for _, tag := range tags {
		tagId, err := strconv.ParseUint(tag, 10, 32)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tagIds = append(tagIds, uint(tagId))
	}
	lastKeySetValue := ctx.Query("lastKeySetValue")
	limit, err := strconv.ParseUint(ctx.Query("limit"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var questions []QuestionDetailResponse
	err = fetchFilteredAndPaginatedQuestions(ctx, uint(tid), sortBy, tagIds, lastKeySetValue, int(limit), &questions)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, questions)
}

func createNewQuestionHandler(ctx *gin.Context) {
	var question *Question

	if err := ctx.ShouldBindJSON(&question); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := middleware.GetUserId(ctx)
	err, userName := users.GetUserNameByID(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var tagIds []uint
	for _, tag := range question.Tags {
		tagIds = append(tagIds, tag.ID)
	}
	err = fetchTagsWithIds(ctx, &question.Tags, tagIds)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	question = &Question{
		UUID:              uuid.New(),
		TopicID:           question.TopicID,
		Title:             question.Title,
		Content:           question.Content,
		CreatedByUser:     userId,
		CreatedByUserName: userName,
		Tags:              question.Tags,
	}

	if err := createQuestion(ctx, question); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// go setQuestionCache(ctx, question)

	ctx.JSON(http.StatusOK, question)
}

func getQuestionWithoutAnswersHandler(ctx *gin.Context) {
	var question Question
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = fetchQuestionWithoutAnswer(ctx, qid, &question)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = stats.UpdateViewHandler(ctx, qid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}
	ctx.JSON(http.StatusOK, question)
}

func getQuestionHandler(ctx *gin.Context) {
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question := &Question{UUID: qid}
	err = fetchQuestionByUUID(ctx, question)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, question)
}

func updateQuestionHandler(ctx *gin.Context) {
	var question Question
	if err := ctx.ShouldBindJSON(&question); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := updateQuestion(ctx, qid, &question); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, question)
}

func updateQuestionTagsHandler(ctx *gin.Context) {
	var requestBody struct {
		Tags []Tag `json:"tags"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tags := requestBody.Tags
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := updateQuestionTags(ctx, qid, &tags); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tags)
}

//	func getQuestionsByTagsHandler(ctx *gin.Context) {
//		qid := ctx.Param("qid")
//		tags := []Tag{}
//		err := fetchQuestionTags(ctx, qid, &tags)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return;
//		}
//		ctx.JSON(http.StatusOK, tags)
//	}
func deleteQuestionHandler(ctx *gin.Context) {
	qid, err := uuid.Parse(ctx.Param("qid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = deleteQuestionByUUID(ctx, qid); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

func getLimitedQuestionByRelevancy(ctx *gin.Context) {
	page, err := strconv.ParseUint(ctx.Param("page"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pageSize, err := strconv.ParseUint(ctx.Param("pageSize"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * pageSize

	var Questions []Question
	err = fetchLimitedQuestionByRelevancy(ctx, int(offset), int(pageSize), &Questions)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Questions)

}
