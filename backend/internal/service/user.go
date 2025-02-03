package service

import (
	"errors"
	"gym/internal/model"
	"gym/pkg/database"
	"gym/pkg/utils"
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

	if !utils.ValidatePassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	return database.DB.Model(&user).Update("password", utils.EncryptPassword(newPassword)).Error
}
