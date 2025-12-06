package app

import (
	"log"
	"team_task_hub/backend/internal/config"
	"team_task_hub/backend/internal/database"
	"team_task_hub/backend/internal/middleware"
	"team_task_hub/backend/internal/repositories"
	"team_task_hub/backend/internal/router"
	"team_task_hub/backend/internal/services"

	_ "team_task_hub/backend/docs"

	"team_task_hub/backend/internal/cache"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	//初始化redis缓存
	if err := cache.InitRedis(cfg); err != nil {
		log.Printf("Redis初始化失败: %v", err)
		log.Println("应用将继续运行，但无缓存功能")
	}

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
	//初始化结构
	userRepo := repositories.NewUserRepository(db)
	emailService := services.NewEmailService(emailConfig, repositories.NewVerificationCodeRepository(db))
	authService := services.NewAuthService(userRepo, emailService, cfg.JWTSecret)

	//接口文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//健康检查路由
	router.SetupRoutes(r, db)
	//邮件路由
	router.SetupEmailRoutes(r.Group("/api"), db, emailConfig)
	//用户认证路由（注册登录，个人信息）
	router.SetupAuthRoutes(r.Group("/api"), db, emailService, cfg.JWTSecret)
	//待办路由
	router.SetupTodoRoutes(r, db, authService)
	//组织路由
	router.SetupOrganizationRoutes(r, db, authService)
}

// Run 启动服务器
func (a *App) Run() {
	//优雅关闭缓存
	defer cache.Close()

	log.Printf("服务器启动在 http://localhost:%s", a.config.ServerPort)
	a.router.Run(":" + a.config.ServerPort)
}
