package monitor_test

import (
	"healthBridge/logger"
	"healthBridge/metrics"
	"healthBridge/monitor"
	"testing"

	"go.uber.org/zap"
)

type MockWebClient struct {
	data string
	err  error
}

func NewMockWebClient(s string, e error) *MockWebClient {
	return &MockWebClient{
		data: s,
		err:  e,
	}
}

func (m *MockWebClient) Get(s string, l *zap.Logger) (string, error) {
	return m.data, m.err
}

func TestMonitorFactory(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	m1, err := monitor.NewMonitor("ping", "mon1", "http://localhost:14271/", 5, log, chan1)

	mockWebClient := NewMockWebClient("data", nil)
	m1.UpdateClient(mockWebClient)

	if err != nil {
		t.Errorf("Factory method failed %v", err)
	}

	name := m1.Identifier()

	if name != "ping_mon1[http://localhost:14271/]" {
		t.Errorf("Got name %s", name)
	}

}

func TestKafdropFactory(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	m1, err := monitor.NewMonitor("kafdrop", "mon1", "http://10.167.197.254:9000/kafdrop/consumer/jaeger-injest-1", 5, log, chan1)
	if err != nil {
		t.Errorf("Factory method failed %v", err)
	}

	name := m1.Identifier()

	if name != "kapdrop_mon1[http://10.167.197.254:9000/kafdrop/consumer/jaeger-injest-1]" {
		t.Errorf("Got name %s", name)
	}
}

func TestUnknownFactory(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	_, err := monitor.NewMonitor("wrongName", "mon", "http://uri/", 5, log, chan1)

	if err == nil {
		t.Errorf("Expected error but factory was created %v", err)
	}
}
