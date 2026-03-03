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

// mockConnConfigurable is a net.Conn whose Read behaviour can be configured
// per test case via readData (success) or readErr (error).
type mockConnConfigurable struct {
	readData []byte
	readErr  error
}

func (c *mockConnConfigurable) Read(b []byte) (int, error) {
	if c.readErr != nil {
		return 0, c.readErr
	}
	return copy(b, c.readData), nil
}
func (c *mockConnConfigurable) Write(b []byte) (int, error)        { return len(b), nil }
func (c *mockConnConfigurable) Close() error                       { return nil }
func (c *mockConnConfigurable) LocalAddr() net.Addr                { return nil }
func (c *mockConnConfigurable) RemoteAddr() net.Addr               { return nil }
func (c *mockConnConfigurable) SetDeadline(t time.Time) error      { return nil }
func (c *mockConnConfigurable) SetReadDeadline(t time.Time) error  { return nil }
func (c *mockConnConfigurable) SetWriteDeadline(t time.Time) error { return nil }

// mockTimeoutError implements net.Error to simulate a read deadline exceeded.
type mockTimeoutError struct{}

func (e *mockTimeoutError) Error() string   { return "i/o timeout" }
func (e *mockTimeoutError) Timeout() bool   { return true }
func (e *mockTimeoutError) Temporary() bool { return true }

// TestSplitPort tests the splitPort function.
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

// TestPopulateMetrics tests the emission of metrics by populateMetrics
func TestPopulateMetrics(t *testing.T) {
	oldNetDialTimeout := netDialTimeout
	defer func() { netDialTimeout = oldNetDialTimeout }()

	args.Network = "tcp"
	args.Address = "localhost:8080"
	args.Timeout = 12

	// populateMetrics shoud set status = 1 when DialTimeout fails
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
	expectedReason := "dial_failed"
	actualReason, ok := ms.Metrics["status_reason"]
	if !ok {
		t.Errorf("populateMetrics, expected status_reason attribute but found none")
	}
	if actualReason != expectedReason {
		t.Errorf("populateMetrics status_reason, got: %s, expected: %s", actualReason, expectedReason)
	}

	// populateMetrics shoud set status = 1 when DialTimeout succeeds
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
	if actual2 != expected2 {
		t.Errorf("populateMetrics, got: %s, expected: %s", actual2, expected2)
	}
	expectedReason = "connected"
	actualReason, ok = ms.Metrics["status_reason"]
	if !ok {
		t.Errorf("populateMetrics, expected status_reason attribute but found none")
	}
	if actualReason != expectedReason {
		t.Errorf("populateMetrics status_reason, got: %s, expected: %s", actualReason, expectedReason)
	}
}

// TestIsUDPNetwork tests the isUDPNetwork function
func TestIsUDPNetwork(t *testing.T) {
	udpNetworks := []string{"udp", "udp4", "udp6"}
	for _, n := range udpNetworks {
		if !isUDPNetwork(n) {
			t.Errorf("isUDPNetwork(%q) should return true", n)
		}
	}

	nonUDPNetworks := []string{"tcp", "tcp4", "tcp6", "ip", "unix", ""}
	for _, n := range nonUDPNetworks {
		if isUDPNetwork(n) {
			t.Errorf("isUDPNetwork(%q) should return false", n)
		}
	}
}

// TestCheckUDPPort tests the checkUDPPort function
func TestCheckUDPPort(t *testing.T) {
	oldNetDialTimeout := netDialTimeout
	defer func() { netDialTimeout = oldNetDialTimeout }()

	// Case 1: dial fails → status=0, reason="dial_failed"
	netDialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return nil, fmt.Errorf("connection refused")
	}
	status, reason := checkUDPPort("udp", "127.0.0.1:9999", 3*time.Second)
	if status != 0 || reason != "dial_failed" {
		t.Errorf("dial failure: got status=%d reason=%q, want 0/dial_failed", status, reason)
	}

	// Case 2: Read returns data → status=1, reason="udp_response_received"
	netDialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return &mockConnConfigurable{readData: []byte("pong")}, nil
	}
	status, reason = checkUDPPort("udp", "127.0.0.1:9999", 3*time.Second)
	if status != 1 || reason != "udp_response_received" {
		t.Errorf("data received: got status=%d reason=%q, want 1/udp_response_received", status, reason)
	}

	// Case 3: Read returns timeout error → status=1, reason="udp_timeout"
	netDialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return &mockConnConfigurable{readErr: &mockTimeoutError{}}, nil
	}
	status, reason = checkUDPPort("udp", "127.0.0.1:9999", 3*time.Second)
	if status != 1 || reason != "udp_timeout" {
		t.Errorf("timeout: got status=%d reason=%q, want 1/udp_timeout", status, reason)
	}

	// Case 4: Read returns non-timeout error (ICMP port unreachable) → status=0, reason="udp_rejected"
	netDialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return &mockConnConfigurable{readErr: fmt.Errorf("connection refused")}, nil
	}
	status, reason = checkUDPPort("udp", "127.0.0.1:9999", 3*time.Second)
	if status != 0 || reason != "udp_rejected" {
		t.Errorf("rejected: got status=%d reason=%q, want 0/udp_rejected", status, reason)
	}
}

// TestPopulateMetric_UDP tests the populateMetrics function for UDP scenarios
func TestPopulateMetrics_UDP(t *testing.T) {
	oldUDPPortChecker := udpPortChecker
	defer func() { udpPortChecker = oldUDPPortChecker }()

	args.Network = "udp"
	args.Address = "127.0.0.1:9999"
	args.Timeout = 3

	cases := []struct {
		name           string
		mockStatus     int
		mockReason     string
		expectedStatus float64
		expectedReason string
	}{
		{"open_filtered", 1, "udp_timeout", 1.0, "udp_timeout"},
		{"response_received", 1, "udp_response_received", 1.0, "udp_response_received"},
		{"rejected", 0, "udp_rejected", 0.0, "udp_rejected"},
		{"dial_failed", 0, "dial_failed", 0.0, "dial_failed"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			udpPortChecker = func(network, address string, timeout time.Duration) (int, string) {
				return tc.mockStatus, tc.mockReason
			}
			ms := metric.NewSet("TestEvent", nil)
			populateMetrics(ms)

			if ms.Metrics["status"] != tc.expectedStatus {
				t.Errorf("status: got %v, want %v", ms.Metrics["status"], tc.expectedStatus)
			}
			if ms.Metrics["status_reason"] != tc.expectedReason {
				t.Errorf("status_reason: got %v, want %v", ms.Metrics["status_reason"], tc.expectedReason)
			}
			if ms.Metrics["network"] != "udp" {
				t.Errorf("network: got %v, want udp", ms.Metrics["network"])
			}
		})
	}
}
