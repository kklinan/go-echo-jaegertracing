package middleware

import (
	"github.com/labstack/echo"

	"github.com/kklinan/go-echo-jaegertracing/jaeger"
)

// Jaeger start span.
func Jaeger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span, _ := jaeger.StartSpanFromHeader(&c.Request().Header, c.Path())
		span.SetTag("api", c.Path())
		defer span.Finish()
		c.Set(jaeger.SpanContextKey, span)
		return next(c)
	}
}
