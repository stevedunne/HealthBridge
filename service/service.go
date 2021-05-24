package service

import (
	"fmt"
	"healthBridge/metrics"
	"healthBridge/monitor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	metrics        *metrics.MetricManager
	signalsChannel chan os.Signal
	monitors       map[string]monitor.Monitor
	adminPort      int
}

// NewService creates a new Service.
func NewService(serverPort int, log *zap.Logger) *Service {
	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, os.Interrupt, syscall.SIGTERM)

	//TODO - initialise metric manager
	metricsManager := metrics.NewMetricManager(log)

	return &Service{
		signalsChannel: signalsChannel,
		adminPort:      serverPort,
		logger:         log,
		metrics:        metricsManager,
		monitors:       make(map[string]monitor.Monitor),
	}
}

func (s *Service) AddMonitor(typeName, name, uri string, pollingInterval int) {
	s.logger.Debug(fmt.Sprintf("Creating new monitor for %s %s %s ", typeName, name, uri))

	m, err := monitor.NewMonitor(typeName, name, uri, pollingInterval, s.logger, s.metrics.Channel)

	if err != nil {
		s.logger.Warn(fmt.Sprintf("Failed to create monitor for %s %s [%s]", typeName, name, uri))
	} else {
		s.monitors[m.Identifier()] = m
		s.logger.Debug(fmt.Sprintf("Added monitor %s [%T] ", m.Identifier(), m))
	}
}

// Start bootstraps the service and starts the admin server.
func (s *Service) Start() error {
	s.logger.Debug("Starting health monitor service")

	//todo - create web servier
	//	   - expose prometheus endpoint from metric manager

	sAdminPort := fmt.Sprintf(":%v", s.adminPort)
	s.logger.Debug(fmt.Sprintf("running /metrics endpoint on %s", sAdminPort))

	http.Handle("/metrics", s.metrics.MetricEndpoint())
	go http.ListenAndServe(sAdminPort, nil)

	s.logger.Debug("Running metrics service")
	go s.metrics.Run()

	// if err := s.Admin.Serve(); err != nil {
	// 	return fmt.Errorf("cannot start the admin server: %w", err)
	// }

	for _, val := range s.monitors {
		s.logger.Debug(fmt.Sprintf("Starting monitor %s [%T] ", val.Identifier(), val))
		val.Start()
	}

	return nil
}

// RunAndThen sets the health check to Ready and blocks
// until SIGTERM is received.
// If then runs the shutdown function and exits.
func (s *Service) RunAndThen(shutdown func()) {

statusLoop:
	for {
		select {
		case <-s.signalsChannel:
			break statusLoop
		}
	}

	log.Print("Shutting down")

	for _, val := range s.monitors {
		val.Destroy()
	}

	if shutdown != nil {
		shutdown()
	}

	log.Print("Shutdown complete")
}
