package monitor

import (
	"healthBridge/metrics"
	"healthBridge/reader"
	"time"

	"go.uber.org/zap"
)

type pingMonitorConfig struct {
	Config
}

func newPingMonitor(name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) *pingMonitorConfig {
	monitor := &pingMonitorConfig{
		Config{
			Name:            name,
			URI:             uri,
			PollingInterval: pollingInterval,
			updateChannel:   ch,
			logger:          log,
			MonitorType:     "ping",
		},
	}
	monitor.WebClient = reader.NewWebClient()
	return monitor
}

func (m *pingMonitorConfig) RunHealthCheck() int {
	m.logger.Debug("Running health check", zap.String("Name", m.Name), zap.String("Uri", m.URI))
	ret := 0

	_, err := m.WebClient.Get(m.URI, m.logger)
	if err != nil {
		m.logger.Warn("Read uri failed", zap.String("Name", m.Name), zap.String("Host", m.getHost()))
		ret = 1
	}

	return ret
}

//The start method will start the internal ticker in the monitor
//kicking off the polling process
//and run a goroutine to handle the tick events and run the health check
func (m *pingMonitorConfig) Start() {
	m.logger.Debug("Starting Monitor", zap.String("Name", m.Name))
	m.ticker = time.NewTicker(time.Second * time.Duration(m.PollingInterval))

	go m.run(m.RunHealthCheck)
}
