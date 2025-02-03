package service

import (
	"encoding/json"
	"fmt"
	"gym/internal/model"
	"gym/pkg/auth"
	"gym/pkg/config"
	"gym/pkg/database"
	"io/ioutil"
	"net/http"
	"time"
)

// WXLoginRequest 微信登录请求
type WXLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录code
}

// WXLoginResponse 微信登录响应
type WXLoginResponse struct {
	Token string      `json:"token"` // JWT token
	User  *model.User `json:"user"`  // 用户信息
}

// code2Session 微信登录凭证校验
type code2Session struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// WechatService 微信服务
type WechatService struct {
	appID     string
	appSecret string
}

// NewWechatService 创建微信服务
func NewWechatService() *WechatService {
	return &WechatService{
		appID:     config.GetConfig().Wechat.AppID,
		appSecret: config.GetConfig().Wechat.AppSecret,
	}
}

// Login 微信小程序登录
func (s *WechatService) Login(req *WXLoginRequest) (*WXLoginResponse, error) {
	// 1. 调用微信登录凭证校验接口
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, req.Code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求微信接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result code2Session
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信接口错误: %s", result.ErrMsg)
	}

	// 2. 查找是否存在用户授权记录
	var userAuth model.UserAuth
	err = database.DB.Where("login_type = ? AND auth_key = ?",
		model.LoginTypeWxMiniApp, result.OpenID).
		Preload("User").First(&userAuth).Error

	// 3. 处理用户信息
	var user model.User
	if err == database.ErrRecordNotFound {
		// 首次登录，创建新用户
		user = model.User{
			Nickname: fmt.Sprintf("微信用户_%s", result.OpenID[:8]), // 默认昵称
			Role:     model.RoleCustomer,
			Status:   1,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}

		// 创建授权记录
		userAuth = model.UserAuth{
			UserID:     user.ID,
			LoginType:  model.LoginTypeWxMiniApp,
			AuthKey:    result.OpenID,
			AuthSecret: result.SessionKey,
			Extra:      fmt.Sprintf(`{"unionid":"%s"}`, result.UnionID),
		}
		if err := database.DB.Create(&userAuth).Error; err != nil {
			return nil, fmt.Errorf("创建授权记录失败: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	} else {
		// 非首次登录，更新session_key和登录时间
		user = userAuth.User
		userAuth.AuthSecret = result.SessionKey
		if result.UnionID != "" {
			userAuth.Extra = fmt.Sprintf(`{"unionid":"%s"}`, result.UnionID)
		}
		if err := database.DB.Save(&userAuth).Error; err != nil {
			return nil, fmt.Errorf("更新授权信息失败: %v", err)
		}

		// 更新最后登录时间
		now := time.Now()
		user.LastLogin = &now
		if err := database.DB.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("更新用户登录时间失败: %v", err)
		}
	}

	// 4. 生成JWT token
	token, err := auth.GenerateToken(user.ID, int(user.Role))
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	return &WXLoginResponse{
		Token: token,
		User:  &user,
	}, nil
}
