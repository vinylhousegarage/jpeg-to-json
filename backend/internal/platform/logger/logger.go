package logger

import "go.uber.org/zap"

func NewLogger(isProduction bool) (*zap.Logger, error) {
    if isProduction {
        return zap.NewProduction()
    }
    return zap.NewDevelopment()
}
