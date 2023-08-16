package posts

import "github.com/gin-gonic/gin"

func fetchAllTopics(ctx *gin.Context, topics *[]Topic) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Find(topics)
	return tx.Error
}

func createTopic(ctx *gin.Context, topic *Topic) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Create(topic)
	return tx.Error
}

func updateTopic(ctx *gin.Context, tid uint, topic *Topic) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Where("id = ?", tid).Updates(Topic{Name: topic.Name})
	return tx.Error
}

func deleteTopic(ctx *gin.Context, tid uint) error {
	tx := db.WithContext(ctx).Model(&Topic{}).Where("id = ?", tid).Delete(&Topic{})
	return tx.Error
}

func fetchAllTopicTags(ctx *gin.Context, tid uint, tags *[]Tag) error {
	tx := db.WithContext(ctx).Model(&Tag{}).Where("topic_id = ?", tid).Find(tags)
	return tx.Error
}

func topicExists(ctx *gin.Context, tid uint) bool {
	var topic Topic
	tx := db.WithContext(ctx).Model(&Topic{}).Where("id = ?", tid).First(&topic)
	return tx.RowsAffected > 0
}