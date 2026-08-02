package logger

import (
	"testing"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	t.Run("Production mode", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{
			App: config.AppConfig{
				Env: "production",
			},
		}

		l, err := NewLogger(cfg)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if l == nil {
			t.Error("Expected logger instance, got nil")
		}
	})

	t.Run("Development mode", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{
			App: config.AppConfig{
				Env: "development",
			},
		}

		l, err := NewLogger(cfg)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if l == nil {
			t.Error("Expected logger instance, got nil")
		}
	})
}
