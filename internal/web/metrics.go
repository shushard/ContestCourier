//nolint:gochecknoglobals // it's just metrics
package web

import (
	"strconv"

	"tester/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const subsystem = "web"

var (
	statusCodeGauge = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: subsystem,
		Name:      "http_response_status_codes_total",
		Help:      "response status codes by path",
	}, []string{"path", "code"})
)

func setStatusCode(path string, code int) {
	statusCodeGauge.WithLabelValues(path, strconv.Itoa(code)).Inc()
}
