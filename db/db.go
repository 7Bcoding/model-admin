package db

import (
	"fmt"
	"llm-ops/models"
	"log"

	"llm-ops/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func InitDB() error {
	cfg := config.Config.MySQL

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库结构
	// err = DB.AutoMigrate(&models.User{})
	// if err != nil {
	// 	return fmt.Errorf("error migrating database: %v", err)
	// }

	log.Printf("Successfully connected to database")

	// 检查是否已存在管理员账户
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)

	// 如果不存在管理员账户，则创建
	if count == 0 {
		adminUser := &models.User{
			Username:    "admin",
			AccountName: "admin",
			Role:        string(models.RoleAdmin),
		}

		if err := adminUser.SetPassword("beta2024"); err != nil {
			return fmt.Errorf("error setting admin password: %v", err)
		}

		if err := DB.Create(adminUser).Error; err != nil {
			return fmt.Errorf("error creating admin user: %v", err)
		}
		log.Printf("Created default admin user with username: admin")
	}

	return nil
}
