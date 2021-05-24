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

type MetricUpdate struct {
	Name string
	Type string
	Host string
	Val  int
}

type MetricManager struct {
	Channel chan MetricUpdate
	stopCh  chan struct{}
}

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

//Provide access to the prometheus http handler to our web server
func (m *MetricManager) MetricEndpoint() http.Handler {
	return promhttp.Handler()
}

//Run the metric manager - it will listen for metrics updates via the channel
func (manager *MetricManager) Run() {
	for {
		select {
		case record := <-manager.Channel:
			logger.Debug(fmt.Sprintf("Received gauge update %s ", record.Name))
			manager.updateGauge(record.Name, record.Type, record.Host, record.Val)
		case <-manager.stopCh:
			return
		}
	}
}

func (manager *MetricManager) Stop() {
	close(manager.stopCh)
}
