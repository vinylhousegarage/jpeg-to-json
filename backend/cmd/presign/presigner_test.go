package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnableCORS(t *testing.T) {
	// 期待するオリジン（テスト用）
	origin := "http://localhost:3000"
    
	// OPTIONSリクエストテスト
	t.Run("OPTIONS request returns 200 and headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		w := httptest.NewRecorder()

		// 関数呼び出し
		enableCORS(w, req)

		// ステータスコード検証
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		// ヘッダー検証
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("expected %s, got %s", origin, got)
		}
	})

	// POSTリクエストテスト
	t.Run("POST request sets headers", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		enableCORS(w, req)

		// ヘッダーがセットされているか確認
		if got := w.Header().Get("Access-Control-Allow-Methods"); got != "POST, OPTIONS" {
			t.Errorf("expected POST, OPTIONS, got %s", got)
		}
	})
}
