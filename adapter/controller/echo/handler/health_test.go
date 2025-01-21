package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"household-account-backend/adapter/controller/echo/handler"
)

func TestHealth(t *testing.T) {
	// Echoインスタンスを作成
	e := echo.New()

	// リクエストとレスポンスのモックを作成
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト対象関数を呼び出し
	err := handler.Health(c)

	// エラーが発生していないことを確認
	assert.NoError(t, err)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// レスポンスボディが期待通りであることを確認
	expected := `{"status":"ok"}`
	assert.JSONEq(t, expected, rec.Body.String())
}
