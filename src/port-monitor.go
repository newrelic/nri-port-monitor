package main

import (
	"net"
	"strings"
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Network string `default:"tcp" help:"Network type. Known networks are tcp, tcp4 (IPv4-only), tcp6 (IPv6-only), udp, udp4 (IPv4-only), udp6 (IPv6-only), ip, ip4 (IPv4-only), ip6 (IPv6-only), unix, unixgram and unixpacket"`
	Address string `default:":80" help:"Address? host:port"`
	Timeout int    `default:"5" help:"Timeout in seconds"`
}

const (
	integrationName    = "com.newrelic.tcp-port-monitor"
	integrationVersion = "0.1.0"
)

var args argumentList

func populateMetrics(ms *metric.MetricSet) error {
	network := strings.TrimSpace(args.Network)
	address := strings.TrimSpace(args.Address)
	status := 0
	conn, err := net.DialTimeout(network, address, time.Duration(args.Timeout)*time.Second)
	if err != nil {
		status = 0
	} else {
		status = 1
		conn.Close()
	}

	ms.SetMetric("network", network, metric.ATTRIBUTE)
	ms.SetMetric("address", address, metric.ATTRIBUTE)
	ms.SetMetric("status", status, metric.GAUGE)
	return nil
}

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	if args.All || args.Metrics {
		ms := integration.NewMetricSet("NetworkPortSample")
		fatalIfErr(populateMetrics(ms))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
