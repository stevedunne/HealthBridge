package monitor

import (
	"healthBridge/metrics"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type kafdropMonitorConfig struct {
	Config
}

// The purpose of this monitor is to extract the
// lag data for a kafka consumer
//
func newKafdropMonitor(name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) *kafdropMonitorConfig {
	monitor := &kafdropMonitorConfig{
		Config{
			Name:            name,
			URI:             uri,
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

	firstNode := nodeList.Get(1)

	firstChild := firstNode.FirstChild
	res := firstChild.Data

	iRes, err := strconv.ParseInt(res, 10, 32)
	if err != nil {
		return -1 //default if the page is not reachable
	}
	return int(iRes)
}

//RunHealthCheck performs the kafdrop health check returning
//the number of unprocessed items or -1 in the event of an error
func (k *kafdropMonitorConfig) RunHealthCheck() int {
	k.logger.Debug("Running health check", zap.String("Name", k.Name), zap.String("Uri", k.URI))

	doc, err := goquery.NewDocument(k.URI)
	if err != nil {
		k.logger.Warn("Failed to read Uri", zap.String("uri", k.URI), zap.Error(err))
		return -1
	}

	return extractLagData(doc)
}

//Start  will start the internal ticker in the monitor
//kicking off the polling process
//and run a goroutine to handle the tick events and run the health check
func (k *kafdropMonitorConfig) Start() {
	k.logger.Debug("Starting Monitor", zap.String("Name", k.Name))
	k.ticker = time.NewTicker(time.Second * time.Duration(k.PollingInterval))

	go k.run(k.RunHealthCheck)
}
