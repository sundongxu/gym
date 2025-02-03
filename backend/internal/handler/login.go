package handler

import (
	"gym/internal/service"

	"github.com/gin-gonic/gin"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	wechatService *service.WechatService
	// TODO: 后续可以添加其他登录服务，如手机号登录服务等
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		wechatService: service.NewWechatService(),
	}
}

// WxMiniLogin 微信小程序登录
func (h *LoginHandler) WxMiniLogin(c *gin.Context) {
	var req service.WXLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	resp, err := h.wechatService.Login(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    resp,
	})
}

// PhoneLogin 手机号密码登录
func (h *LoginHandler) PhoneLogin(c *gin.Context) {
	// TODO: 实现手机号密码登录
	c.JSON(200, gin.H{
		"code":    200,
		"message": "该功能尚未实现",
	})
}

// SmsLogin 短信验证码登录
func (h *LoginHandler) SmsLogin(c *gin.Context) {
	// TODO: 实现短信验证码登录
	c.JSON(200, gin.H{
		"code":    200,
		"message": "该功能尚未实现",
	})
}
