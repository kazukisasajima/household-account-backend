package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"household-account-backend/pkg/logger"
)

// カラーコード
var (
	green  = "\033[32m"
	blue   = "\033[34m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

// CustomRequestLogger はリクエスト情報をログ出力するミドルウェアです。
func CustomRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			status := c.Response().Status
			method := c.Request().Method

			// ログの出力
			logger.ZapLogger.Info("Request",
				zap.String("method", method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", status),
				zap.String("latency", stop.Sub(start).String()),
				zap.String("client_ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			)

			// カラフルな出力
			fmt.Printf("[%s%s%s] %s%3d%s | %13v | %15s | %s%-7s%s %s\n",
				blue, time.Now().Format("2006/01/02 - 15:04:05"), reset,
				getStatusColor(status), status, reset,
				stop.Sub(start),
				c.RealIP(),
				getMethodColor(method), method, reset,
				c.Request().URL.Path,
			)

			return err
		}
	}
}

// CustomRecovery はパニックをログに記録するミドルウェアです。
func CustomRecovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.ZapLogger.Error("Panic recovered",
						zap.Any("error", r),
					)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// ステータスコードに応じた色を返す関数
func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return blue
	case status >= 400 && status < 500:
		return yellow
	default:
		return red
	}
}

// HTTPメソッドに応じた色を返す関数
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return green
	case "PUT":
		return yellow
	case "DELETE":
		return red
	default:
		return reset
	}
}