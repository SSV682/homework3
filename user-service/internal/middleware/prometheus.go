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
			Name:      "response_time_seconds",
			Help:      "Duration of HTTP requests",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	errorRate = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "echo",
			Subsystem: "http",
			Name:      "error_rate_by_user_service",
			Help:      "Count of 5xx",
		},
	)
)

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(errorRate)
}

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
			if statusCode == 500 {
				errorRate.Inc()
			}
			totalRequests.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
			return err
		}

	}
}
