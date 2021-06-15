package monitor_test

import (
	"errors"
	"healthBridge/logger"
	"healthBridge/metrics"
	"healthBridge/monitor"
	"testing"
)

func TestPingFactory_success(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	m1, err := monitor.NewMonitor("ping", "mon_1", "http://localhost:14271/", 5, log, chan1)

	mockWebClient := NewMockWebClient("data", nil)
	m1.UpdateClient(mockWebClient)

	if err != nil {
		t.Errorf("Factory method failed %v", err)
	}

	val := m1.RunHealthCheck()

	if val != 0 {
		t.Errorf("Expected zero but got %v", val)
	}
}

func TestPingFactory_fail(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	m1, err := monitor.NewMonitor("ping", "mon_1", "http://localhost:14271/", 5, log, chan1)

	mockWebClient := NewMockWebClient("data", errors.New("testfailure"))
	m1.UpdateClient(mockWebClient)

	if err != nil {
		t.Errorf("Factory method failed %v", err)
	}

	val := m1.RunHealthCheck()

	if val != 1 {
		t.Errorf("Expected failure but got %v", val)
	}
}
