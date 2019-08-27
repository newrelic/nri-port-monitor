package main

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

type mockConn struct{}

func (c *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}
func (c *mockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}
func (c *mockConn) Close() error {
	return nil
}
func (c *mockConn) LocalAddr() net.Addr {
	return nil
}
func (c *mockConn) RemoteAddr() net.Addr {
	return nil
}
func (c *mockConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestSplitPort(t *testing.T) {
	// Test without a port
	expected := "80"
	actual := splitPort("mangatsika")
	if actual != expected {
		t.Errorf("splitPort, got: %s, expected: %s", actual, expected)
	}

	// Test with a valid port
	expected = "8080"
	actual = splitPort("mangatsika:8080")
	if actual != expected {
		t.Errorf("splitPort, got: %s, expected: %s", actual, expected)
	}
}

func TestPopulateMetrics(t *testing.T) {
	oldNetDialTimeout := netDialTimeout
	defer func() { netDialTimeout = oldNetDialTimeout }()

	args.Network = "tcp"
	args.Address = "localhost:8080"
	args.Timeout = 12

	// populateMetrics shoud set status = 1 when DialTimeout fails.
	netDialTimeout = func(
		network, address string,
		timeout time.Duration,
	) (net.Conn, error) {
		return nil, fmt.Errorf("mangatsika")
	}
	ms := metric.NewSet("TestEvent", nil)

	expected := 0.0
	populateMetrics(ms)
	actual, ok := ms.Metrics["status"]
	if !ok {
		t.Errorf("populateMetrics, expected status attribute but found none")
	}
	if actual != expected {
		t.Errorf("populateMetrics, got: %f, expected: %f", actual, expected)
	}

	// populateMetrics shoud set status = 1 when DialTimeout succeeds.
	netDialTimeout = func(
		network, address string,
		timeout time.Duration,
	) (net.Conn, error) {
		return &mockConn{}, nil
	}
	ms = metric.NewSet("TestEvent", nil)

	expected = 1.0
	populateMetrics(ms)
	actual, ok = ms.Metrics["status"]
	if !ok {
		t.Errorf("populateMetrics, expected status attribute but found none")
	}
	if actual != expected {
		t.Errorf("populateMetrics, got: %f, expected: %f", actual, expected)
	}

	expected2 := "tcp"
	actual2, ok := ms.Metrics["network"]
	if !ok {
		t.Errorf("populateMetrics, expected network attribute but found none")
	}
	if actual2 != expected2 {
		t.Errorf("populateMetrics, got: %s, expected: %s", actual2, expected2)
	}

	expected2 = "localhost:8080"
	actual2, ok = ms.Metrics["address"]
	if !ok {
		t.Errorf("populateMetrics, expected address attribute but found none")
	}
	if actual2 != expected2 {
		t.Errorf("populateMetrics, got: %s, expected: %s", actual2, expected2)
	}

	expected2 = "8080"
	populateMetrics(ms)
	actual2, ok = ms.Metrics["port"]
	if !ok {
		t.Errorf("populateMetrics, expected port attribute but found none")
	}
	if actual != expected {
		t.Errorf("populateMetrics, got: %f, expected: %f", actual, expected)
	}
}
