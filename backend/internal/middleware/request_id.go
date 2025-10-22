package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	RqeuestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(RqeuestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			c.Set(RequestIDKey, requestID)
			c.Response().Header().Set(RqeuestIDHeader, requestID)

			return next(c)
		}
	}
}

func GetRequestID(c echo.Context) string {
	if requestID, ok := c.Get(RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}
