package monitoring

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	Hits     *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

func RegisterMonitoring(server *echo.Echo) *PrometheusMetrics {
	var metrics = new(PrometheusMetrics)

	metrics.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	metrics.Duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "duration",
	}, []string{"status", "path", "method"})

	prometheus.MustRegister(metrics.Hits, metrics.Duration)

	server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return metrics
}
