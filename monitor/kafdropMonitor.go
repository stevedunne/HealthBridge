package monitor

import (
	"healthBridge/metrics"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type kafdropMonitorConfig struct {
	MonitorConfig
}

// The purpose of this monitor is to extract the
// lag data for a kafka consumer
//
func newKafdropMonitor(name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) *kafdropMonitorConfig {
	monitor := &kafdropMonitorConfig{
		MonitorConfig{
			Name:            name,
			Uri:             uri,
			PollingInterval: pollingInterval,
			updateChannel:   ch,
			logger:          log,
			MonitorType:     "kapdrop",
		},
	}
	return monitor
}

func extractLagData(doc *goquery.Document) int {

	nodeList := doc.Find("#topic-0-table td b")
	res := nodeList.Get(1).FirstChild.Data

	iRes, err := strconv.ParseInt(res, 10, 32)
	if err != nil {
		return -1 //default if the page is not reachable
	}
	return int(iRes)
}

func (m *kafdropMonitorConfig) RunHealthCheck() int {
	m.logger.Debug("Running health check", zap.String("Name", m.Name), zap.String("Uri", m.Uri))

	doc, err := goquery.NewDocument(m.Uri)
	if err != nil {
		m.logger.Warn("Failed to read Uri", zap.String("uri", m.Uri), zap.Error(err))
		return -1
	}

	return extractLagData(doc)
}

//The start method will start the internal ticker in the monitor
//kicking off the polling process
//and run a goroutine to handle the tick events and run the health check
func (k *kafdropMonitorConfig) Start() {
	k.logger.Debug("Starting Monitor", zap.String("Name", k.Name))
	k.ticker = time.NewTicker(time.Second * time.Duration(k.PollingInterval))

	go k.run(k.RunHealthCheck)
}
