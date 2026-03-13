package models

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleOperator UserRole = "operator"
	RoleUser     UserRole = "user"
	RoleGuest    UserRole = "guest"
)

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Username     string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	AccountName  string    `gorm:"type:varchar(255)" json:"account_name"`
	Password     string    `gorm:"-" json:"password"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) CheckPassword(password string) bool {
	// use md5
	hash := md5.Sum([]byte(password))
	return u.PasswordHash == hex.EncodeToString(hash[:])
}

func (u *User) SetPassword(password string) error {
	// use md5
	hash := md5.Sum([]byte(password))
	u.PasswordHash = hex.EncodeToString(hash[:])
	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}

func (r UserRole) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleOperator:
		return "operator"
	case RoleUser:
		return "user"
	default:
		return "guest"
	}
}

// IsValidRole 检查角色是否有效
func IsValidRole(role string) bool {
	validRoles := []string{"admin", "operator", "user", "guest"}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}
