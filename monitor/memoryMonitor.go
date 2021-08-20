//the purpose of the memory monitor is to watch the metrics on a go ap
//and in the event of a spike in memory consumption it will save a dump of the heap diagnostics
//data to disk
//this requires the go application to expose the prometheus metrics endpoint and the golang pprof endpoint

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

type memoryMonitorConfig struct {
	Config
}

//factory method
func newMemoryMonitor(name, uri string, pollingInterval int, log *zap.Logger, ch chan<- metrics.MetricUpdate) *memoryMonitorConfig {
	monitor := &memoryMonitorConfig{
		Config{
			Name:            name,
			URI:             uri,
			PollingInterval: pollingInterval,
			updateChannel:   ch,
			logger:          log,
			MonitorType:     "memoryMonitor",
		},
	}
	monitor.WebClient = reader.NewWebClient()
	return monitor
}

//Run health check method

//RunHealthCheck performs the kafdrop health check returning 0 when no issue
//is detected or 1 in the event of a spike. A spike will also result in a dump of the heap
//being saved to disk
func (m *memoryMonitorConfig) RunHealthCheck() int {
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
			//	break - dont break as we actaully want the last one
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
	m.logger.Info("Virtual Memory for agent: ", zap.String("Name", m.Name), zap.Float64("memoryUsage", memoryUsage), zap.String("Host", m.getHost()))
	if memoryUsage > 500000000 {
		//save the file

		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		rootPath := path + "\\Data\\" + getFolderName(m.getHost())

		err = os.Mkdir(rootPath, 0755)
		if err != nil {
			log.Println(err)
		}
		m.readUriToFile(m.URI+"/debug/pprof/heap", rootPath+"\\heap")
		m.readUriToFile(m.URI+"/debug/pprof/allocs", rootPath+"\\allocs")
		m.readUriToFile(m.URI+"/version", rootPath+"\\version.json")
		m.readUriToFile(m.URI, rootPath+"\\health.json")

		return 1 //the counter will reflect how many dump files we have gotten
	}
	return 0
}

//Start  will start the internal ticker in the monitor
//kicking off the polling process
//and run a goroutine to handle the tick events and run the health check
func (m *memoryMonitorConfig) Start() {
	m.logger.Debug("Starting Monitor", zap.String("Name", m.Name))
	m.ticker = time.NewTicker(time.Second * time.Duration(m.PollingInterval))

	go m.run(m.RunHealthCheck)
}

func (m *memoryMonitorConfig) readUriToFile(uri, filename string) {
	data, err := m.WebClient.Get(uri, m.logger)
	if err != nil {
		m.logger.Warn("Read data failed: ", zap.String("Uri", uri), zap.String("filename", filename), zap.String("error", err.Error()))
	}

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		m.logger.Warn("Error writing data to file", zap.String("Name", filename), zap.String("error", err.Error()))
	}

}

func getFolderName(host string) string {
	t := time.Now()

	str := fmt.Sprintf("%s-%s", host, t.Format(time.RFC3339))
	str = strings.Replace(str, ":", "_", -1)
	return str
}
