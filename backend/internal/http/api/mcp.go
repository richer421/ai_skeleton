package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/richer/ai_skeleton/internal/common"
	"github.com/richer/ai_skeleton/internal/mcp"
)

var mcpAdapter mcp.MCPAdapter

// InitMCP 初始化 MCP 适配器（在 router setup 时调用一次）
func InitMCP(adapter mcp.MCPAdapter) {
	mcpAdapter = adapter
}

// MCPListTools 列出所有 MCP 工具
// @Summary 列出 MCP 工具
// @Description 列出所有已注册的 MCP 工具
// @Tags MCP
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=[]mcp.ToolSchema}
// @Router /api/v1/mcp/tools [get]
func MCPListTools(c *gin.Context) {
	tools := mcpAdapter.ListTools()
	c.JSON(http.StatusOK, common.Success(tools))
}

// MCPExecute 执行 MCP 工具
// @Summary 执行 MCP 工具
// @Description 执行指定的 MCP 工具
// @Tags MCP
// @Accept json
// @Produce json
// @Param request body mcp.MCPRequest true "MCP 请求"
// @Success 200 {object} common.Response{data=mcp.MCPResponse}
// @Router /api/v1/mcp/execute [post]
func MCPExecute(c *gin.Context) {
	var req mcp.MCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(400, "invalid request: "+err.Error()))
		return
	}

	result, err := mcpAdapter.HandleRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}
