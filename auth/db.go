package auth

import (
	"fmt"
	"net/http"

	"github.com/campus-fora/middleware"

	"github.com/campus-fora/constants"
	"github.com/gin-gonic/gin"
)

type UpdateRoleRequest struct {
	UserID    string         `json:"user_id" binding:"required"`
	NewRoleID constants.Role `json:"new_role_id" binding:"required"`
}

func fetchUser(ctx *gin.Context, user *User, userID string) error {
	tx := DB.WithContext(ctx).Where("ID = ?", userID).First(&user)
	return tx.Error
}
func getUserRole(ctx *gin.Context, userID string) (constants.Role, error) {
	var user User
	tx := DB.WithContext(ctx).Where("ID = ?", userID).First(&user)
	return user.Role, tx.Error
}
func updatePassword(ctx *gin.Context, userID string, password string) (bool, error) {
	tx := DB.WithContext(ctx).Model(&User{}).Where("ID = ?", userID).Update("password", password)
	return tx.RowsAffected > 0, tx.Error
}
func UserExists(user_id string) bool {
	var user User
	result := DB.First(&user, "ID = ?", fmt.Sprint(user_id))
	return result.Error == nil

}
func toggleBlock(ctx *gin.Context, ID uint) (bool, error) {
	var currStatus bool
	tx := DB.WithContext(ctx).Model(&User{}).Where("ID = ?", ID).Select("blocked").First(&currStatus)
	if tx.Error != nil {
		return false, tx.Error
	}
	// fmt.Printf(currentStatus.First())
	tx = DB.WithContext(ctx).Model(&User{}).Where("ID = ?", ID).Update("blocked", !currStatus)
	return !currStatus, tx.Error
}
func getRoleAndStatus(ctx *gin.Context, userID string) (constants.Role, bool, error) {
	var user User
	tx := DB.WithContext(ctx).Where("ID = ? AND Blocked = ?", userID, true).First(&user)
	return user.Role, user.Blocked, tx.Error
}
func updateRoleByAdmin(ctx *gin.Context, ID string, roleID constants.Role) error {
	tx := DB.WithContext(ctx).Model(&User{}).Where("ID = ?", ID).Update("role_id", roleID)
	return tx.Error
}

func updateUserRole(ctx *gin.Context) {

	var updateReq UpdateRoleRequest

	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentRoleID constants.Role
	currentRoleID, err := getUserRole(ctx, updateReq.UserID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	middleware.Authenticator()(ctx)
	middleware.EnsureAuthority()(ctx)
	if middleware.GetUserID(ctx) == "" {
		return
	}
	var user_id = middleware.GetUserID(ctx)

	userRole, _, err := getRoleAndStatus(ctx, user_id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if userRole > currentRoleID || userRole > updateReq.NewRoleID || userRole > 101 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this user's role"})
		return
	}

	err = updateRoleByAdmin(ctx, updateReq.UserID, updateReq.NewRoleID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}
