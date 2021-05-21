package main

import (
	"fmt"

	"healthBridge/logger"
	"healthBridge/service"
	"healthBridge/ymlReader"
)

/*
# ToDo:
#  - Read config from yml
#  - create service to contain all relevant processes
#  - Read 'ping' end point for success/failure
#  - create prometheus end point
#  - run pingers as golang routines
#  	- use channel to send back data
#  - manage params cobra/viper
#  - logging - set level
*/

func main() {

	//TODO - validate args

	// fmt.Println("Started")
	// fmt.Printf("Got args %v\r\n", len(os.Args))
	// if len(os.Args) == 0 {
	// 	panic("No path supplied for config file")
	// }
	// for i, val := range os.Args {
	// 	fmt.Printf("%v: %s\r\n", i, val)
	// }
	// var configLocation = os.Args[0]
	// fmt.Printf("Got file path %s\r\n", configLocation)

	//	fmt.Println(ymlReader.WriteConf())

	//TODO - initialise logger

	logger, err := logger.NewLogger()
	if err != nil {
		panic("Could not initialize logger")
	}

	configLocation := "C:\\src\\HealthBridge\\config.yml"
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
