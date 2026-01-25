package mcp

import (
	"context"
	"log"

	"github.com/richer/ai_skeleton/internal/service/health"
)

// RegisterAllTools 统一注册所有 MCP 工具
func RegisterAllTools(adapter MCPAdapter) error {
	// 注册健康检查工具
	if err := registerHealthTool(adapter); err != nil {
		log.Printf("Failed to register health tool: %v", err)
		return err
	}

	// 在这里添加更多工具注册
	// if err := registerUserTool(adapter); err != nil {
	//     log.Printf("Failed to register user tool: %v", err)
	//     return err
	// }

	log.Printf("Successfully registered %d MCP tools", len(adapter.ListTools()))
	return nil
}

// registerHealthTool 注册健康检查工具
func registerHealthTool(adapter MCPAdapter) error {
	schema := ToolSchema{
		Name:        "health_check",
		Description: "检查系统健康状态",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
			"required":   []string{},
		},
	}

	handler := func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// 每次调用时创建 service
		svc := health.NewHealthService()
		return svc.Check(ctx)
	}

	return adapter.RegisterTool("health_check", schema, handler)
}
