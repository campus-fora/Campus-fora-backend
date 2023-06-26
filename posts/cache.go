package posts

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type thread_as_hmap struct {
	thread_detail []byte `redis:"thread_detail"`
	posts         []byte `redis:"posts"`
	tags          []byte `redis:"tags"`
}

func getAllQuestionDetailsCache(ctx *gin.Context, threads *[]all_thread_response) error {
	cursor := uint64(0)
	size, err := rdb.DBSize(ctx).Result()
	if err != nil {
		return err
	}
	if size == 0 {
		return redis.Nil
	}

	for {
		keys, nextCursor, err := rdb.Scan(ctx, cursor, "thread:*", 10).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			var thread_detail thread_detail
			var tags []Tag

			cmd, err := rdb.HGet(ctx, key, "thread_detail").Bytes()
			if err != nil {
				return err
			}

			b := bytes.NewReader(cmd)

			if err := gob.NewDecoder(b).Decode(&thread_detail); err != nil && err.Error() != "EOF" {
				return err
			}

			cmd, err = rdb.HGet(ctx, key, "tags").Bytes()
			if err != nil {
				return err
			}

			b = bytes.NewReader(cmd)

			if err := gob.NewDecoder(b).Decode(&tags); err != nil && err.Error() != "EOF" {
				return err
			}

			thread_detail_with_tags := all_thread_response{
				ID:            thread_detail.ID,
				CreatedAt:     thread_detail.CreatedAt,
				Title:         thread_detail.Title,
				Content:       thread_detail.Content,
				CreatedByUser: thread_detail.CreatedByUser,
				Tags:          tags,
			}
			*threads = append(*threads, thread_detail_with_tags)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	fmt.Println("cache hit")
	return nil
}

func setQuestionCache(ctx *gin.Context, thread Thread) error {
	var b_thread_detail, b_Posts, b_Tags bytes.Buffer

	thread_detail := thread_detail{
		ID:            thread.ID,
		CreatedAt:     thread.CreatedAt,
		Title:         thread.Title,
		Content:       thread.Content,
		CreatedByUser: thread.CreatedByUser,
	}

	err := gob.NewEncoder(&b_thread_detail).Encode(thread_detail)
	if err != nil {
		fmt.Print("error in cahcing: ", err.Error())
		return err
	}

	err = gob.NewEncoder(&b_Posts).Encode(thread.Posts)
	if err != nil {
		return err
	}

	err = gob.NewEncoder(&b_thread_detail).Encode(thread.Tags)
	if err != nil {
		fmt.Print("error in cahcing: ", err.Error())
		return err
	}

	err = rdb.HSet(ctx, "thread:"+strconv.FormatUint(uint64(thread.ID), 10),
		"thread_detail", b_thread_detail.Bytes(),
		"posts", b_Posts.Bytes(),
		"tags", b_Tags.Bytes()).Err()
	if err != nil {
		fmt.Print("error in cahcing: ", err.Error())
	}
	fmt.Print("Caching completed")
	return err
}

func setQuestionDetailsCache(ctx *gin.Context, thread Thread) error {
	var b_thread_detail, b_Tags bytes.Buffer

	thread_detail := thread_detail{
		ID:            thread.ID,
		CreatedAt:     thread.CreatedAt,
		Title:         thread.Title,
		Content:       thread.Content,
		CreatedByUser: thread.CreatedByUser,
	}

	err := gob.NewEncoder(&b_thread_detail).Encode(thread_detail)
	if err != nil {
		return err
	}

	err = gob.NewEncoder(&b_thread_detail).Encode(thread.Tags)
	if err != nil {
		return err
	}

	return rdb.HSet(ctx, "thread:"+strconv.FormatUint(uint64(thread.ID), 10),
		"thread_detail", b_thread_detail.Bytes(),
		"tags", b_Tags.Bytes()).Err()
}

func setAllQuestionDetailCache(ctx *gin.Context, threads []Thread) error {
	for _, thread := range threads {
		err := setQuestionDetailsCache(ctx, thread)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("caching complete")
	return nil
}
