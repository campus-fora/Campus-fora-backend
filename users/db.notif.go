package users

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SaveNewToken(ctx *gin.Context, userID uint, token string, deviceId string) error {
	MAX_DEVICES := viper.GetInt64("NOTIFICATION.MAX_DEVICES")
	err := Db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Association("NotifTokens").Append(&NotifTokens{Token: token, DeviceId: deviceId})
	if err != nil {
		return err
	}
	cnt := Db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Association("NotifTokens").Count()
	if cnt > MAX_DEVICES {
		go deleteOldestToken(ctx, userID)
	}
	return nil
}

func DeleteToken(ctx *gin.Context, userID uint, deviceId string) error {
	err := Db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Association("NotifTokens").Delete(&NotifTokens{}, "device_id = ?", deviceId)
	if err != nil {
		return err
	}
	return Db.WithContext(ctx).Unscoped().Model(&NotifTokens{}).Where("device_id = ?", deviceId).Delete(&NotifTokens{}).Error
}

func deleteOldestToken(ctx *gin.Context, userID uint) error {
	var user User
	var token NotifTokens
	err := Db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Preload("NotifTokens", func() {
		Db.WithContext(ctx).Order("created_at asc").Limit(1)
	}).First(&user).Error
	if err != nil {
		return err
	}
	if len(user.NotifTokens) > 0 {
		token = user.NotifTokens[0]
	} else {
		return nil
	}
	return Db.WithContext(ctx).Unscoped().Model(&NotifTokens{}).Where("token = ?", token.Token).Delete(&NotifTokens{}).Error
}
