package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gym/internal/model"
	"gym/pkg/database"
)

// UserService 用户服务
type UserService struct{}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile 更新用户信息
func (s *UserService) UpdateUserProfile(userID uint, nickname, phone, avatar string) error {
	updates := map[string]interface{}{
		"nickname": nickname,
		"phone":    phone,
		"avatar":   avatar,
	}
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if !s.validatePassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	return database.DB.Model(&user).Update("password", s.encryptPassword(newPassword)).Error
}

// encryptPassword 密码加密
func (s *UserService) encryptPassword(password string) string {
	// 使用SHA256进行密码加密
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// validatePassword 验证密码
func (s *UserService) validatePassword(password, hashedPassword string) bool {
	return s.encryptPassword(password) == hashedPassword
}
