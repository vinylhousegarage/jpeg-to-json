package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	t.Run("Production mode", func(t *testing.T) {
		t.Parallel()

		l, err := NewLogger(true)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if l == nil {
			t.Error("Expected logger instance, got nil")
		}
	})

	t.Run("Development mode", func(t *testing.T) {
		t.Parallel()

		l, err := NewLogger(false)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if l == nil {
			t.Error("Expected logger instance, got nil")
		}
	})
}
