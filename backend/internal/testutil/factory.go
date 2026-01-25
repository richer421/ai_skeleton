package testutil

import (
	"github.com/richer/ai_skeleton/internal/common"
)

// CreateTestHealthResponse 创建测试用健康检查响应
func CreateTestHealthResponse(opts ...func(*common.HealthResponse)) *common.HealthResponse {
	resp := &common.HealthResponse{
		Status:    "ok",
		Timestamp: "2026-01-25T10:00:00Z",
		Version:   "1.0.0",
	}

	for _, opt := range opts {
		opt(resp)
	}

	return resp
}

// WithStatus 设置状态
func WithStatus(status string) func(*common.HealthResponse) {
	return func(r *common.HealthResponse) {
		r.Status = status
	}
}

// WithVersion 设置版本
func WithVersion(version string) func(*common.HealthResponse) {
	return func(r *common.HealthResponse) {
		r.Version = version
	}
}
