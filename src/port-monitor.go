package main

import (
	"net"
	"strings"
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Network string `default:"tcp" help:"Network type. Known networks are tcp, tcp4 (IPv4-only), tcp6 (IPv6-only), udp, udp4 (IPv4-only), udp6 (IPv6-only), ip, ip4 (IPv4-only), ip6 (IPv6-only), unix, unixgram and unixpacket"`
	Address string `default:":80" help:"Address in the format host:port"`
	Timeout int    `default:"5" help:"Timeout in seconds"`
}

const (
	integrationName    = "com.newrelic.labs.nri-port-monitor"
	integrationVersion = "1.4.1"
)

var (
	args           argumentList
	netDialTimeout = net.DialTimeout
	udpPortChecker = checkUDPPort
)

// splitPort splits an address into host and port, returning the port or "80" if not specified.
func splitPort(address string) string {
	slices := strings.Split(address, ":")
	if len(slices) == 1 {
		return "80"
	}

	return slices[1]
}

// isUDPNetwork returns true when the network string refers to a UDP variant.
func isUDPNetwork(network string) bool {
	switch network {
	case "udp", "udp4", "udp6":
		return true
	}
	return false
}

// checkUDPPort probes a UDP port by sending a small packet and waiting for
// either a response or an ICMP port-unreachable error surfaced as a read
// error.
//
// Returns the status (1 = open/filtered, 0 = closed) and a reason string.
//
// Response map:
//   - data received             → 1, "udp_response_received"
//   - read deadline exceeded    → 1, "udp_open"     (open or silently filtered; matches nmap open|filtered)
//   - non-timeout read error    → 0, "udp_rejected" (ICMP port unreachable)
//   - dial or write error       → 0, "dial_failed"
func checkUDPPort(network, address string, timeout time.Duration) (int, string) {
	conn, err := netDialTimeout(network, address, timeout)
	if err != nil {
		return 0, "dial_failed"
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	if _, err := conn.Write([]byte{0x00}); err != nil {
		return 0, "dial_failed"
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// No ICMP rejection within the window — port is open or filtered (likely via firewall).
			return 1, "udp_open"
		}
		// Non-timeout error (e.g. ECONNREFUSED from ICMP port unreachable).
		return 0, "udp_rejected"
	}

	return 1, "udp_response_received"
}

func populateMetrics(ms *metric.Set) {
	network := strings.TrimSpace(args.Network)
	address := strings.TrimSpace(args.Address)
	port := splitPort(address)
	status := 0
	statusReason := ""
	timeout := time.Duration(args.Timeout) * time.Second

	if isUDPNetwork(network) {
		status, statusReason = udpPortChecker(network, address, timeout)
	} else {
		conn, err := netDialTimeout(network, address, timeout)
		if err == nil {
			status = 1
			statusReason = "connected"
			conn.Close()
		} else {
			statusReason = "dial_failed"
		}
	}

	ms.SetMetric("network", network, metric.ATTRIBUTE)
	ms.SetMetric("address", address, metric.ATTRIBUTE)
	ms.SetMetric("port", port, metric.ATTRIBUTE)
	ms.SetMetric("status", status, metric.GAUGE)
	ms.SetMetric("status_reason", statusReason, metric.ATTRIBUTE)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	integration, err := integration.New(
		integrationName,
		integrationVersion,
		integration.Args(&args),
	)
	panicOnErr(err)

	entity := integration.LocalEntity()

	args.NriAddHostname = true
	if args.All() || args.Metrics {
		populateMetrics(entity.NewMetricSet("NetworkPortSample"))
	}

	panicOnErr(integration.Publish())
}
