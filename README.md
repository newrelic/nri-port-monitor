[![Community Project header](https://github.com/newrelic/open-source-office/raw/master/examples/categories/images/Community_Project.png)](https://github.com/newrelic/open-source-office/blob/master/examples/categories/index.md#community-project)

# New Relic Infrastructure On-Host Integration for monitoring Network Ports

Reports up or down status for a network (TCP, UDP etc) port.

## Requirements
- Infrastructure agent installed - see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)

## Installation

1. Download a [pre-generated binary](https://github.com/newrelic/nri-port-monitor/releases) for your OS or build from source (see [Building](#building))
2. Copy the integration binary (`nri-port-monitor`) to `/var/db/newrelic-infra/custom-integrations/`
3. Copy the [sample configuration](port-monitor-config.yml) to `/etc/newrelic-infra/integrations.d/`
4. Set required environment variables (see [Configuration](#configuration))
5. Restart the Infrastructure Agent:
   ```bash
   sudo systemctl restart newrelic-infra
   ```

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NETWORK` | Yes | `tcp` | Network type (e.g., tcp, udp, unix) |
| `ADDRESS` | Yes | `localhost:80` | Target address in host:port format |
| `TIMEOUT` | No | `5` | Connection timeout in seconds |


> **NOTE** If using an integration v1.3 or lower, use [legacy config files](./legacy/). These should be copied as instructed in step 3 within the Installation section.

## View data

By issuing the following NRQL, you can display the results of the port monitor.

```sql NRQL
SELECT latest(status), latest(status_reason) FROM NetworkPortSample FACET address SINCE 30 MINUTES AGO TIMESERIES
```

### `status` values
0 = Port closed
1 = Port open

### `status_reason` values

| `status_reason`         | `status` | Description |
|-------------------------|----------|-------------|
| `connected`             | 1        | TCP connection succeeded |
| `dial_failed`           | 0        | Could not dial or write to the target |
| `udp_response_received` | 1        | UDP — data was received from the target |
| `udp_open`              | 1        | UDP — no response within the deadline; port is open or silently filtered (likely via firewall) |
| `udp_rejected`          | 0        | UDP — ICMP port unreachable received |


## Building

### Prerequisites

- Go 1.23 or higher
- Make

### Steps

1. Clone this repo and open a terminal in the folder where you cloned it.
2. Run `make` to build for all platforms, or:
   - `make macos-arm` for macOS ARM64
   - `make macos-intel` for macOS AMD64
   - `make linux-arm` for Linux ARM64
   - `make linux-intel` for Linux AMD64
   - `make windows` for Windows AMD64

The generated binaries will be located in the `bin/` directory.

To test the integration, run `nri-port-monitor` with passed configuration env variables:

```bash
$ NETWORK=tcp ADDRESS=newrelic.com:443 TIMEOUT=5  ./bin/<build>/nri-port-monitor
```

The output should be json similar to:

```json
{"name":"com.newrelic.tcp-port-monitor","protocol_version":"3","integration_version":"3.0.0","data":[{"metrics":[{"address":"newrelic.com:443","event_type":"NetworkPortSample","network":"tcp","port":"443","status":1}],"inventory":{},"events":[]}]}
```

## Support

<a href="https://github.com/newrelic?q=nrlabs-viz&amp;type=all&amp;language=&amp;sort="><img src="https://user-images.githubusercontent.com/1786630/214122263-7a5795f6-f4e3-4aa0-b3f5-2f27aff16098.png" height=50 /></a>

This project is actively maintained by the New Relic Labs team. Connect with us directly by [creating issues](../../issues) or [asking questions in the discussions section](../../discussions) of this repo.

We also encourage you to bring your experiences and questions to the [Explorers Hub](https://discuss.newrelic.com) where our community members collaborate on solutions and new ideas.

New Relic has open-sourced this project, which is provided AS-IS WITHOUT WARRANTY OR DEDICATED SUPPORT.

## Security

As noted in our [security policy](https://github.com/newrelic/nr-labs-pages/security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in one of our products or websites, we welcome and greatly appreciate you reporting it to New Relic through our [our bug bounty program](https://docs.newrelic.com/docs/security/security-privacy/information-security/report-security-vulnerabilities/).

## Contributing

Contributions are welcome (and if you submit a Enhancement Request, expect to be invited to contribute it yourself :grin:). Please review our [Contributors Guide](CONTRIBUTING.md).

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## License

This project is distributed under the [Apache 2 license](LICENSE).
