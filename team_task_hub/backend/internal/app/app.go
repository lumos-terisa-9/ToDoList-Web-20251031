package app

import (
	"log"
	"team_task_hub/backend/internal/config"
	"team_task_hub/backend/internal/database"
	"team_task_hub/backend/internal/middleware"
	"team_task_hub/backend/internal/repositories"
	"team_task_hub/backend/internal/router"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App 应用主体
type App struct {
	router *gin.Engine
	config *config.Config
}

// New 创建应用实例
func New() *App {
	log.Println("启动 Team Task Hub 服务...")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db := initDatabase(cfg)

	// 设置Gin模式
	setGinMode(cfg)

	// 创建路由
	r := gin.Default()
	r.Use(middleware.CORS())

	// 设置路由
	setupRoutes(r, db, cfg)

	return &App{
		router: r,
		config: cfg,
	}
}

// initDatabase 初始化数据库
func initDatabase(cfg *config.Config) *gorm.DB {
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	return db
}

// setGinMode 设置Gin模式
func setGinMode(cfg *config.Config) {
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
		log.Println("运行在调试模式")
	}
}

// setupRoutes 设置路由
func setupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// 创建邮箱配置
	emailConfig := &services.EmailConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
	}

	// 设置路由
	emailService := services.NewEmailService(emailConfig, repositories.NewVerificationCodeRepository(db))
	router.SetupRoutes(r, db)
	router.SetupEmailRoutes(r.Group("/api"), db, emailConfig)
	router.SetupAuthRoutes(r.Group("/api"), db, emailService, cfg.JWTSecret)
}

// Run 启动服务器
func (a *App) Run() {
	log.Printf("服务器启动在 http://localhost:%s", a.config.ServerPort)
	a.router.Run(":" + a.config.ServerPort)
}
