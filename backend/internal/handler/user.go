package handler

import (
	"gym/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户控制器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户控制器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: &service.UserService{},
	}
}

// GetProfile 获取用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    1003,
			"message": "获取用户信息失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    user,
	})
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	type UpdateRequest struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.userService.UpdateUserProfile(userID, req.Nickname, req.Phone, req.Avatar); err != nil {
		c.JSON(500, gin.H{
			"code":    1004,
			"message": "更新用户信息失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(500, gin.H{
			"code":    1005,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "修改成功",
	})
}
