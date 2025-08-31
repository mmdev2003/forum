package middleware

import (
	"forum-authentication/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"time"
)

func LoggerMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(request echo.Context) error {
			start := time.Now()

			requestID := request.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}
			request.Request().Header.Set("X-Request-ID", requestID)

			logger := slog.With(
				slog.String("request_id", requestID),
				slog.String("method", request.Request().Method),
				slog.String("path", request.Request().URL.Path),
				slog.String("service_name", "forum-dialog"),
			)
			request.Set(model.LoggerKey, logger)

			logger.Info("request started")
			err := next(request)
			duration := float64(time.Since(start)) / 1e6

			if err != nil {
				logger.Error("request failed",
					slog.Any("error", err),
					slog.Float64("duration_ms", duration),
				)
			} else {
				logger.Info("request completed",
					slog.Int("status", request.Response().Status),
					slog.Float64("duration_ms", duration),
				)
			}
			return err
		}
	}
}
