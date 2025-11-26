package main

import (
	"team_task_hub/backend/internal/app"
)

// @title Team Task Hub API
// @version 1.0
// @description 团队任务管理平台API
// @host localhost:8080
// @BasePath
func main() {
	app.New().Run()
}
