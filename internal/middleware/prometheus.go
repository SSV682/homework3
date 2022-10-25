package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "echo",
			Subsystem: "http",
			Name:      "http_request_total",
			Help:      "Number of get requests",
		},
		[]string{"method", "path", "status"},
	)

	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "echo",
			Subsystem: "http",
			Name:      "http_response_time_seconds",
			Help:      "Duration of HTTP requests",
			Buckets: []float64{
				0.0005,
				0.001, // 1ms
				0.002,
				0.005,
				0.01, // 10ms
				0.02,
				0.05,
				0.1, // 100 ms
				0.2,
				0.5,
				1.0, // 1s
				2.0,
				5.0,
				10.0, // 10s
				15.0,
				20.0,
				30.0,
			},
		},
		[]string{"method", "path"},
	)
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			method := ctx.Request().Method
			path := ctx.Path()
			timer := prometheus.NewTimer(httpDuration.WithLabelValues(method, path))
			err := next(ctx)
			timer.ObserveDuration()

			if err != nil {
				ctx.Error(err)
			}
			statusCode := ctx.Response().Status
			totalRequests.WithLabelValues(method, path, strconv.Itoa(statusCode))
			return err
		}

	}
}
