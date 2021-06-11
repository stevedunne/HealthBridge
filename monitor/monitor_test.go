package monitor

import (
	"healthBridge/logger"
	"healthBridge/metrics"
	"testing"
)

func TestPingFactory(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	m1, err := NewMonitor("ping", "mon_1", "http://localhost:14271/", 5, log, chan1)

	if err != nil {
		t.Errorf("Factory method failed %v", err)
	}

	name := m1.Identifier()

	if name != "mon_1" {
		t.Errorf("Got name %s", name)
	}

	val := m1.RunHealthCheck()

	if val == -1 {
		t.Errorf("Expected positive but got %v", val)
	}
}

func TestKafdropFactory(t *testing.T) {

	log, _ := logger.NewLogger()
	chan1 := make(chan metrics.MetricUpdate)

	m1, err := NewMonitor("kafdrop", "mon_1", "http://10.167.197.254:9000/kafdrop/consumer/jaeger-injest-1", 5, log, chan1)
	if err != nil {
		t.Errorf("Factory method failed %v", err)
	}

	name := m1.Identifier()

	if name != "mon_1" {
		t.Errorf("Got name %s", name)
	}

	val := m1.RunHealthCheck()

	if val == -1 {
		t.Errorf("Expected positive but got %v", val)
	} else {
		t.Logf("Value returned was %v", val)
	}

}
