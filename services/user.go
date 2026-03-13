package services

import (
	"llm-ops/models"
	"log"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() []models.User {
	var users []models.User
	result := s.DB.Select("id, username, account_name, role, created_at, updated_at").Find(&users)
	if result.Error != nil {
		log.Printf("Error fetching users: %v", result.Error)
		return []models.User{}
	}
	return users
}

// AddUser 添加新用户
func (s *UserService) AddUser(user *models.User) error {
	return s.DB.Create(user).Error
}

// UpdateUserRole 更新用户角色
func (s *UserService) UpdateUserRole(username string, role string) error {
	return s.DB.Model(&models.User{}).Where("username = ?", username).Update("role", role).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(username string) error {
	return s.DB.Where("username = ?", username).Delete(&models.User{}).Error
}

// GetUserByUsername 通过用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("username = ? OR account_name = ?", username, username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UserExists(username string) (bool, error) {
	var user models.User
	err := s.DB.Where("username = ? OR account_name = ?", username, username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err == nil {
		return true, nil
	}
	return false, nil
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(userID int, newPassword string) error {
	// 获取用户
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// 更新数据库
	return s.DB.Model(&user).Update("password_hash", user.PasswordHash).Error
}
