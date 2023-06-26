package posts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type all_thread_response struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedByUser string    `json:"created_by_user"`
	Tags          []Tag     `json:"tags"`
}

func getAllThreadsDetail(ctx *gin.Context) {
	var threads []all_thread_response

	err := getAllQuestionDetailsCache(ctx, &threads)
	if err == nil {
		fmt.Print("cache hit")
		ctx.JSON(http.StatusOK, threads)
		return
	}
	if err == redis.Nil {
		var threads []Thread
		err = fetchAllThreadDetails(ctx, &threads)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		go setAllQuestionDetailCache(ctx, threads)
		fmt.Println("get request completed")
		ctx.JSON(http.StatusOK, threads)
		return
	}
	fmt.Print(err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func postThread(ctx *gin.Context) {
	var thread Thread

	if err := ctx.ShouldBindJSON(&thread); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		print(err.Error())
		return
	}

	if err := createThread(ctx, &thread); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go setQuestionCache(ctx, thread)
	fmt.Println("post request completed")

	ctx.JSON(http.StatusOK, thread)
}

func getPosts(ctx *gin.Context) {
	response := struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}{
		Question: "What is the meaning of life?",
		Answer:   "42",
	}

	ctx.JSON(http.StatusOK, response)
}
