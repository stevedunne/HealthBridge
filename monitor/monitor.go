package monitor

import (
	"errors"
	"fmt"
	"healthBridge/metrics"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type Monitor interface {
	Start()
	Destroy()
	RunHealthCheck() int
	Identifier() string
}

type MonitorConfig struct {
	Name            string
	Uri             string
	PollingInterval int
	MonitorType     string

	logger        *zap.Logger
	ticker        *time.Ticker // periodic ticker
	updateChannel chan<- metrics.MetricUpdate

	Monitor
}

//factory method to create the specified type of monitor
func NewMonitor(monitorType, name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) (Monitor, error) {
	//logger.Debug()
	switch monitorType {
	case "ping":
		return newPingMonitor(name, uri, pollingInterval, log, ch), nil
	case "kafdrop":
		return newKafdropMonitor(name, uri, pollingInterval, log, ch), nil
	default:
		return nil, errors.New("No such monitor type '%s'")
	}
}

// Stops the ticker
func (m *MonitorConfig) Destroy() {
	m.logger.Debug("Destroying Monitor", zap.String("Name", m.Name))
	m.ticker.Stop()
}

// Stops the ticker
func (m *MonitorConfig) Identifier() string {
	return fmt.Sprintf("%s_%s[%s]", m.MonitorType, m.Name, m.Uri)
}

func (m *MonitorConfig) getHost() string {
	url, err := url.Parse(m.Uri)
	if err != nil {
		m.logger.Debug("Could not parse Uri", zap.Error(err))
		return m.Uri
	}
	return url.Host
}

//run is an internal method that will run as a goroutine running the
//health check and publishing the results to the channel
func (m *MonitorConfig) run(healthCheck func() int) {
	m.logger.Debug("Running Monitor", zap.String("Name", m.Name))
	host := m.getHost()
	for {
		select {
		case <-m.ticker.C:
			m.logger.Debug("Tick recieved for monitor", zap.String("Name", m.Name))
			m.logger.Debug(fmt.Sprintf("Got monitor %T", m))
			i := healthCheck()
			m.updateChannel <- metrics.MetricUpdate{Name: m.Name, Host: host, Type: m.MonitorType, Val: i}
		}
	}
}
