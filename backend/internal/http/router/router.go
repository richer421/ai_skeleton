package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/richer/ai_skeleton/internal/http/api"
	"github.com/richer/ai_skeleton/internal/http/middleware"
	"github.com/richer/ai_skeleton/internal/mcp"
)

// Setup 设置路由
func Setup() *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 初始化 MCP 适配器并注册所有工具
	mcpAdapter := mcp.NewMCPAdapter()
	if err := mcp.RegisterAllTools(mcpAdapter); err != nil {
		log.Fatalf("Failed to register MCP tools: %v", err)
	}
	api.InitMCP(mcpAdapter)

	// API 路由组
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", api.HealthCheck)

		// MCP 协议
		mcpGroup := v1.Group("/mcp")
		{
			mcpGroup.GET("/tools", api.MCPListTools)
			mcpGroup.POST("/execute", api.MCPExecute)
		}
	}

	return r
}
