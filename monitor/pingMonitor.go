package monitor

import (
	"healthBridge/metrics"
	"net/http"
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
	return monitor
}

func (m *pingMonitorConfig) RunHealthCheck() int {
	m.logger.Debug("Running health check", zap.String("Name", m.Name), zap.String("Uri", m.URI))
	ret := 0
	//resp, err := http.Get(m.Uri)
	resp, err := http.Get(m.URI)
	if err != nil {
		m.logger.Warn("Read uri failed: ", zap.String("error", err.Error()))
		ret = 1
	}

	if resp != nil && resp.Body != nil {
		resp.Body.Close()
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
