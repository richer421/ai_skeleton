package health

import (
	"context"
	"time"

	"github.com/richer/ai_skeleton/internal/common"
)

// HealthService 健康检查服务接口
type HealthService interface {
	Check(ctx context.Context) (*common.HealthResponse, error)
}

type healthService struct {
	version string
}

// NewHealthService 创建健康检查服务
func NewHealthService() HealthService {
	return &healthService{
		version: "1.0.0",
	}
}

// Check 执行健康检查
func (s *healthService) Check(ctx context.Context) (*common.HealthResponse, error) {
	return &common.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   s.version,
	}, nil
}
