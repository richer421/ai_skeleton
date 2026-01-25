package health

import (
	"context"
	"testing"

	"github.com/richer/ai_skeleton/internal/testutil"
)

func TestHealthService_Check(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "健康检查成功",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewHealthService()
			result, err := s.Check(context.Background())

			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
				testutil.AssertNotNil(t, result)
				testutil.AssertEqual(t, result.Status, "ok")
				testutil.AssertEqual(t, result.Version, "1.0.0")
			}
		})
	}
}
