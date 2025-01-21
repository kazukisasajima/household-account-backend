package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func TimeoutMiddleware(duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), duration)
			defer cancel()

			c.SetRequest(c.Request().WithContext(ctx))

			errChan := make(chan error, 1)
			go func() {
				errChan <- next(c)
			}()

			select {
			case err := <-errChan:
				return err
			case <-ctx.Done():
				return c.JSON(http.StatusRequestTimeout, map[string]string{
					"message": "timeout",
				})
			}
		}
	}
}