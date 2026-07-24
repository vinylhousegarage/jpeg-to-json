package logger

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.AppEnv == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
