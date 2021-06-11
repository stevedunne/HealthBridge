package main

import (
	"flag"
	"fmt"
	"os"

	"healthBridge/logger"
	"healthBridge/service"
	"healthBridge/ymlreader"

	_ "net/http/pprof"
)

func main() {

	var flagvar int
	flag.IntVar(&flagvar, "port", 8080, "the port number for the server to run on. Default 8080.")

	//	configLocation := "C:\\src\\HealthBridge\\config.yml"
	var configLocation string
	flag.StringVar(&configLocation, "config", "", "the path to the yml config file describing the monitors")

	logger.Flags()
	flag.Parse()

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
	//var configLocation = os.Args[1]
	//	fmt.Println(ymlreader.WriteConf())
	//	configLocation := "C:\\src\\HealthBridge\\config.yml"

	fmt.Printf("Got file path %s\r\n", configLocation)

	if len(configLocation) == 0 {
		panic("No path supplied for config file")
	}

	conf, err := ymlreader.ReadConf(configLocation)
	if err != nil {
		panic(fmt.Errorf("Error loading yml configuration %v", err))
	}

	fmt.Printf("Created config %+v\r\n", conf)

	//Todo - create service(logger, monitorList )
	svc := service.NewService(8080, logger)
	for _, val := range conf.HealthConfigs {

		for _, uri := range val.MonitorType.URI {
			svc.AddMonitor(val.MonitorType.Type, val.MonitorName, uri, int(conf.GlobalConfig.ScrapeInterval))
		}
	}

	svc.Start()

	svc.RunAndThen(func() { fmt.Println("... stopping ... ") })

}
