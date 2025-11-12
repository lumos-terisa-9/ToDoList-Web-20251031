package main

import (
	"log"
	"team_task_hub/backend/internal/config"
	"team_task_hub/backend/internal/database"
	"team_task_hub/backend/internal/middleware"
	"team_task_hub/backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("启动 Team Task Hub 后端服务...")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置Gin模式
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
		log.Println("运行在调试模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin实例
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 设置路由（这里包含了 /health 路由）
	router.SetupRoutes(r, db)

	// 启动服务器
	log.Printf("🌐 服务器启动在 http://localhost:%s", cfg.ServerPort)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
