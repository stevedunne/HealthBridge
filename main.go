package main

import (
	"fmt"
	"os"

	"healthBridge/logger"
	"healthBridge/service"
	"healthBridge/ymlReader"
)

func main() {

	//TODO - initialise logger

	logger, err := logger.NewLogger()
	if err != nil {
		panic("Could not initialize logger")
	}

	//TODO - validate args

	logger.Debug("Starting Health Monitor")
	if len(os.Args) == 0 {
		panic("No path supplied for config file")
	}
	for i, val := range os.Args {
		logger.Debug(fmt.Sprintf("%v: %s", i, val))
	}
	var configLocation = os.Args[1]
	fmt.Printf("Got file path %s\r\n", configLocation)

	//	fmt.Println(ymlReader.WriteConf())
	//configLocation := "C:\\src\\HealthBridge\\config.yml"

	if len(configLocation) == 0 {
		panic("No path supplied for config file")
	}

	conf, err := ymlReader.ReadConf(configLocation)
	if err != nil {
		panic(fmt.Errorf("Error loading yml configuration %v", err))
	}

	fmt.Printf("Created config %+v\r\n", conf)

	//Todo - create service(logger, monitorList )
	svc := service.NewService(8080, logger)
	for _, val := range conf.HealthConfigs {

		for _, uri := range val.MonitorType.Uri {
			svc.AddMonitor(val.MonitorType.Type, val.MonitorName, uri, int(conf.GlobalConfig.ScrapeInterval))
		}
	}

	svc.Start()

	svc.RunAndThen(func() { fmt.Println("... stopping ... ") })

}
