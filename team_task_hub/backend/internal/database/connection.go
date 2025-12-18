package database

import (
	"fmt"
	"log"
	"team_task_hub/backend/internal/config"
	"team_task_hub/backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.Todo{},
		&models.TodoAssignment{},
		&models.Activity{},
		&models.ActivityParticipation{},
		&models.ScoreRecord{},
		&models.Permission{},
		&models.VerificationCode{},
		&models.OrganizationApplication{},
	)
	if err != nil {
		return nil, fmt.Errorf("表迁移失败: %v", err)
	}

	log.Println("数据库连接和表迁移成功!")
	return db, nil
}
