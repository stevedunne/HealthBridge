package monitor

import (
	"errors"
	"fmt"
	"healthBridge/metrics"
	"healthBridge/reader"
	"net/url"

	"time"

	"go.uber.org/zap"
)

//Monitor provides the key methods to be implemented by any healy monitor
type Monitor interface {
	Start()
	Destroy()
	RunHealthCheck() int
	Identifier() string

	UpdateClient(reader.IWebClient)
}

//Config represents the config details for an individial monitor
type Config struct {
	Name            string
	URI             string
	PollingInterval int
	MonitorType     string

	WebClient reader.IWebClient

	logger        *zap.Logger
	ticker        *time.Ticker // periodic ticker
	updateChannel chan<- metrics.MetricUpdate
}

//NewMonitor factory method to create a specified type of monitor
func NewMonitor(monitorType, name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) (Monitor, error) {
	switch monitorType {
	case "ping":
		return newPingMonitor(name, uri, pollingInterval, log, ch), nil
	case "kafdrop":
		return newKafdropMonitor(name, uri, pollingInterval, log, ch), nil
	case "jaeger-agent":
		return newJaegerAgentMonitor(name, uri, 120, log, ch), nil
	default:
		return nil, errors.New("No such monitor type '%s'")
	}
}

//Destroy stops the ticker
func (m *Config) Destroy() {
	m.logger.Debug("Destroying Monitor", zap.String("Name", m.Name))
	m.ticker.Stop()
}

// Identifier creates an id for the monitor from its config
func (m *Config) Identifier() string {
	return fmt.Sprintf("%s_%s[%s]", m.MonitorType, m.Name, m.URI)
}

func (m *Config) getHost() string {
	url, err := url.Parse(m.URI)
	if err != nil {
		m.logger.Debug("Could not parse Uri", zap.Error(err))
		return m.URI
	}
	return url.Host
}

//run is an internal method that will run as a goroutine running the
//health check and publishing the results to the channel
func (m *Config) run(healthCheck func() int) {
	m.logger.Debug("Running Monitor", zap.String("Name", m.Name))
	host := m.getHost()
	for {
		select {
		case <-m.ticker.C:
			m.logger.Debug("Tick received for monitor", zap.String("Name", m.Name))
			m.logger.Debug(fmt.Sprintf("Got monitor %T", m))
			i := healthCheck()
			m.updateChannel <- metrics.MetricUpdate{Name: m.Name, Host: host, Type: m.MonitorType, Val: i}
		}
	}
}

//UpdateClient sets an alternate webClient - use for testing
func (m *Config) UpdateClient(client reader.IWebClient) {
	m.logger.Debug("Updating web client")
	m.WebClient = client
}
