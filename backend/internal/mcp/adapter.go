package mcp

import (
	"context"
	"encoding/json"
)

// MCPAdapter MCP 协议适配器接口
type MCPAdapter interface {
	// RegisterTool 注册 MCP 工具
	RegisterTool(name string, schema ToolSchema, handler ToolHandler) error

	// HandleRequest 处理 MCP 请求
	HandleRequest(ctx context.Context, req *MCPRequest) (*MCPResponse, error)

	// ListTools 列出所有已注册的工具
	ListTools() []ToolSchema
}

// ToolSchema MCP 工具描述
type ToolSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolHandler 工具处理函数
type ToolHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// MCPRequest MCP 请求
type MCPRequest struct {
	Tool   string                 `json:"tool"`
	Params map[string]interface{} `json:"params"`
}

// MCPResponse MCP 响应
type MCPResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// mcpAdapter MCP 适配器实现
type mcpAdapter struct {
	tools    map[string]ToolSchema
	handlers map[string]ToolHandler
}

// NewMCPAdapter 创建 MCP 适配器
func NewMCPAdapter() MCPAdapter {
	return &mcpAdapter{
		tools:    make(map[string]ToolSchema),
		handlers: make(map[string]ToolHandler),
	}
}

// RegisterTool 注册工具
func (a *mcpAdapter) RegisterTool(name string, schema ToolSchema, handler ToolHandler) error {
	a.tools[name] = schema
	a.handlers[name] = handler
	return nil
}

// HandleRequest 处理请求
func (a *mcpAdapter) HandleRequest(ctx context.Context, req *MCPRequest) (*MCPResponse, error) {
	handler, exists := a.handlers[req.Tool]
	if !exists {
		return &MCPResponse{
			Success: false,
			Error:   "tool not found",
		}, nil
	}

	result, err := handler(ctx, req.Params)
	if err != nil {
		return &MCPResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &MCPResponse{
		Success: true,
		Data:    result,
	}, nil
}

// ListTools 列出所有工具
func (a *mcpAdapter) ListTools() []ToolSchema {
	tools := make([]ToolSchema, 0, len(a.tools))
	for _, schema := range a.tools {
		tools = append(tools, schema)
	}
	return tools
}

// ToJSONSchema 将 ToolSchema 转换为 JSON Schema 格式
func (t *ToolSchema) ToJSONSchema() (string, error) {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
