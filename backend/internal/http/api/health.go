package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/richer/ai_skeleton/internal/common"
	"github.com/richer/ai_skeleton/internal/service/health"
)

// HealthCheck 健康检查接口
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=common.HealthResponse}
// @Router /api/v1/health [get]
func HealthCheck(c *gin.Context) {
	svc := health.NewHealthService()
	result, err := svc.Check(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}
