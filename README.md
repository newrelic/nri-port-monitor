# New Relic Infrastructure Integration for monitoring Network Ports

Reports up or down status for a network (TCP, UDP etc) port

## Requirements

### New Relic Infrastructure Agent

This is the description about how to run the Port Monitor Integration with New Relic Infrastructure agent, so it is required to have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

## Installation

* Download an archive file for the port-monitor Integration
* Place the executables under `bin` directory and the definition file `port-monitor-definition.yml` in `/var/db/newrelic-infra/newrelic-integrations`
* Set execution permissions for the binary file `nr-port-monitor`
* Place the integration configuration file `port-monitor-config.yml.sample` in `/etc/newrelic-infra/integrations.d`

## Configuration

In order to use the Port Monitor Integration it is required to configure `port-monitor-config.yml.sample` file. Firstly, rename the file to `port-monitor-config.yml`. Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

You can view your data in Insights by creating your own custom NRQL queries. To
do so use **NetworkPortSample** event types.

## Compatibility

* Supported OS: linux
* Edition: