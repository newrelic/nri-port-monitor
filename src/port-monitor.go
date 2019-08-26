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
	Address string `default:":80" help:"Address? host:port"`
	Timeout int    `default:"5" help:"Timeout in seconds"`
}

const (
	integrationName    = "com.newrelic.tcp-port-monitor"
	integrationVersion = "2.0.0"
)

var (
	args           argumentList
	netDialTimeout = net.DialTimeout
)

func splitPort(address string) string {
	slices := strings.Split(address, ":")
	if len(slices) == 1 {
		return "80"
	}

	return slices[1]
}

func populateMetrics(ms *metric.Set) {
	network := strings.TrimSpace(args.Network)
	address := strings.TrimSpace(args.Address)
	port := splitPort(address)
	status := 0

	conn, err := netDialTimeout(
		network,
		address,
		time.Duration(args.Timeout)*time.Second,
	)

	if err == nil {
		status = 1
		conn.Close()
	}

	ms.SetMetric("network", network, metric.ATTRIBUTE)
	ms.SetMetric("address", address, metric.ATTRIBUTE)
	ms.SetMetric("port", port, metric.ATTRIBUTE)
	ms.SetMetric("status", status, metric.GAUGE)
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
