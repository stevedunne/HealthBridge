//the purpose of the jaeger agent monitor is to watch the metrics on a jeager agent
//and in the event of a spike in memory consumption it will save a dump of the heap diagnostics
//data to disk

package monitor

import (
	"fmt"
	"healthBridge/metrics"
	"healthBridge/reader"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type jaegerAgentMonitorConfig struct {
	Config
}

//factory method
func newJaegerAgentMonitor(name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) *jaegerAgentMonitorConfig {
	monitor := &jaegerAgentMonitorConfig{
		Config{
			Name:            name,
			URI:             uri,
			PollingInterval: pollingInterval,
			updateChannel:   ch,
			logger:          log,
			MonitorType:     "jaegerAgent",
		},
	}
	monitor.WebClient = reader.NewWebClient()
	return monitor
}

//Run health check method

//RunHealthCheck performs the kafdrop health check returning 0 when no issue
//is detected or 1 in the event of a spike. A spike will also result in a dump of the heap
//being saved to disk
func (m *jaegerAgentMonitorConfig) RunHealthCheck() int {
	m.logger.Debug("Running health check", zap.String("Name", m.Name), zap.String("Uri", m.URI))

	//read the metrics end point [http://localhost:14271/metrics]
	uri := m.URI + "/metrics"
	metricData, err := m.WebClient.Get(uri, m.logger)
	if err != nil {
		m.logger.Warn("Read uri failed: ", zap.String("Name", m.Name), zap.String("error", err.Error()))
		return 0
	}

	//extract the memory data
	/*
	   # HELP process_resident_memory_bytes Resident memory size in bytes.
	   # TYPE process_resident_memory_bytes gauge
	   process_resident_memory_bytes 2.0742144e+07
	   # HELP process_virtual_memory_bytes Virtual memory size in bytes.
	   # TYPE process_virtual_memory_bytes gauge
	   process_virtual_memory_bytes 7.38127872e+08
	   # HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
	   # TYPE process_virtual_memory_max_bytes gauge
	   process_virtual_memory_max_bytes 1.8446744073709552e+19
	*/

	strMem := ""
	lines := strings.Split(metricData, "\n")
	for _, v := range lines {
		if strings.Contains(v, "process_virtual_memory_bytes") {
			strMem = v
			break
		}
	}

	if len(strMem) == 0 {
		m.logger.Warn("Virtual Memory data not found", zap.String("Name", m.Name))
		return 0
	}

	pair := strings.Split(strMem, " ")
	memoryUsage, err := strconv.ParseFloat(pair[1], 64)
	if err != nil {
		m.logger.Warn("Could not extract Virtual Memory data", zap.String("Name", m.Name), zap.String("strMem", strMem))
		return 0
	}

	//if its above a set threshold then
	//	get the heap diag info  [http://localhost:14271/debug/pprof/heap?debug=1]
	//  save to disk
	// return a 1 if a file is saved
	if memoryUsage > 1000000000 {
		//save the file

		//read the metrics end point [http://localhost:14271/metrics]
		uri := m.URI + "/debug/pprof/heap?debug=1"
		perfData, err := m.WebClient.Get(uri, m.logger)
		if err != nil {
			m.logger.Warn("Read perf data failed: ", zap.String("Name", m.Name), zap.String("error", err.Error()))
			return 0
		}
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		path = path + "\\Data\\" + getFileName(m.getHost())
		err = ioutil.WriteFile(path, []byte(perfData), 0644)
		if err != nil {
			m.logger.Warn("Error writing perf data to file", zap.String("Name", m.Name), zap.String("error", err.Error()))
			return 0
		}
		return 1 //the counter will reflect how many dump files we have gotten
	}
	return 0
}

//Start  will start the internal ticker in the monitor
//kicking off the polling process
//and run a goroutine to handle the tick events and run the health check
func (m *jaegerAgentMonitorConfig) Start() {
	m.logger.Debug("Starting Monitor", zap.String("Name", m.Name))
	m.ticker = time.NewTicker(time.Second * time.Duration(m.PollingInterval))

	go m.run(m.RunHealthCheck)
}

func getFileName(host string) string {
	t := time.Now()

	str := fmt.Sprintf("%s-%s.log", host, t.Format(time.RFC3339))
	str = strings.Replace(str, ":", "_", -1)
	return str
}
