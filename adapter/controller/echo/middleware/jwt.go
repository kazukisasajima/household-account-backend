package middleware

import (
	"fmt"
	"net/http"
	"os"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"household-account-backend/pkg/logger"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cookieから"auth_token"を取得
			// クライアントから送信されたリクエスト内のCookieを調べ、"auth_token"を取得する。
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Missing auth_token cookie"})
			}

			// JWTトークンを解析して署名を検証
			// トークンの署名が正しいかどうかを確認し、有効性を検証する。
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token"})
			}

			// トークンのClaimsを型変換し、正しい形式（jwt.MapClaims）であることを確認
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				logger.Info("Parsed JWT Token Claims: " + fmt.Sprintf("%v", claims))
				// 後続の処理で利用できるようにトークン全体をコンテキストに保存
				c.Set("user", token)
			} else {
				logger.Error("Invalid token claims")
				return c.JSON(http.StatusUnauthorized, "Invalid Claims")
			}

			return next(c)
		}
	}
}