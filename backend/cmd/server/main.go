package main

import (
	"fmt"
	"log"

	"github.com/richer/ai_skeleton/internal/config"
	"github.com/richer/ai_skeleton/internal/http/router"
)

// @title AI Skeleton API
// @version 1.0
// @description AI 全栈脚手架 API 文档
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置路由
	r := router.Setup()

	// 启动服务器
	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
