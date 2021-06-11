package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	logger      *zap.Logger
	healthGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "healthbridge_gauge",
		Help: "Health gauge help.",
	}, []string{"Name", "Type", "Host"})
	gaugeHitCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "healthbridge_hit_count",
		Help: "healthbridge health test hit counter.",
	}, []string{"Name", "Type", "Host"})
)

// MetricUpdate - this encapsulates the information for a specifi cmetric
type MetricUpdate struct {
	Name string
	Type string
	Host string
	Val  int
}

// MetricManager contains the channel to receive metric data
type MetricManager struct {
	Channel chan MetricUpdate
	stopCh  chan struct{}
}

// NewMetricManager creates a new MetricManager instance
func NewMetricManager(l *zap.Logger) *MetricManager {
	logger = l

	mm := &MetricManager{
		Channel: make(chan MetricUpdate, 1024),
		stopCh:  make(chan struct{}),
	}

	err := prometheus.Register(healthGauge)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to register health gauge  %v", err))
	} else {
		logger.Debug("registered health gauge")
	}

	err = prometheus.Register(gaugeHitCount)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to register gauge hit counter  %v", err))
	} else {
		logger.Debug("registered gauge hit counter")
	}

	return mm
}

//find a named gauge in the collection and update its value
func (m *MetricManager) updateGauge(name, metricType, host string, gaugeVal int) {
	logger.Debug(fmt.Sprintf("Updating gauges %s %v", name, gaugeVal))

	healthGauge.WithLabelValues(name, metricType, host).Set(float64(gaugeVal))

	gaugeHitCount.WithLabelValues(name, metricType, host).Inc()
}

// MetricEndpoint provides access to the prometheus http handler to our web server
func (m *MetricManager) MetricEndpoint() http.Handler {
	return promhttp.Handler()
}

//Run the metric manager - it will listen for metrics updates via the channel
func (m *MetricManager) Run() {
	for {
		select {
		case record := <-m.Channel:
			logger.Debug(fmt.Sprintf("Received gauge update %s ", record.Name))
			m.updateGauge(record.Name, record.Type, record.Host, record.Val)
		case <-m.stopCh:
			return
		}
	}
}

// Stop closes the metricManagers stop channel
func (m *MetricManager) Stop() {
	close(m.stopCh)
}
