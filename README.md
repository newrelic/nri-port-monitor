[![Community Project header](https://github.com/newrelic/open-source-office/raw/master/examples/categories/images/Community_Project.png)](https://github.com/newrelic/open-source-office/blob/master/examples/categories/index.md#community-project)

# New Relic Infrastructure On-Host Integration for monitoring Network Ports

Reports up or down status for a network (TCP, UDP etc) port.

## Requirements

You should have the infrastructure agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

## Installation

* Download and unpack the ZIP file from [Releases](https://github.com/newrelic/nri-port-monitor/releases).  

```sh bash
wget https://github.com/newrelic/nri-port-monitor/releases/download/1.3/nri-port-monitor.tar.gz
tar -zxvf nri-port-monitor.tar.gz 
```

* Copy the `bin` directory with `nri-port-monitor` executable, and the `port-monitor-definition.yml` config file to `/var/db/newrelic-infra/newrelic-integrations`.  

```sh bash
sudo cp nri-port-monitor/bin/port-monitor /var/db/newrelic-infra/newrelic-integrations/bin/
sudo cp nri-port-monitor/port-monitor-definition.yml /var/db/newrelic-infra/newrelic-integrations/
```

* Set execution permissions for the binary file `nr-port-monitor`.  

``` sh bash
sudo chmod +x /var/db/newrelic-infra/newrelic-integrations/bin/port-monitor
```

* Place the integration configuration file `port-monitor-config.yml.sample` in `/etc/newrelic-infra/integrations.d`.

## Configuration

In order to use the Port Monitor Integration it is required to configure `port-monitor-config.yml.sample` file. Firstly, rename the file to `port-monitor-config.yml`.  

```sh bash
sudo cp nri-port-monitor/port-monitor-config.yml.sample /etc/newrelic-infra/integrations.d/port-monitor-config.yml
```

Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

```sh bash
sudo systemctl restart newrelic-infra.service
```

Data should start flowing into your New Relic account. See [Understand and use data from Infrastructure integrations](https://docs.newrelic.com/docs/integrations/infrastructure-integrations/get-started/understand-use-data-infrastructure-integrations).

## View data

By issuing the following NRQL, you can display the results of the port monitor.

```sql NRQL
SELECT latest(status) FROM NetworkPortSample FACET address SINCE 30 MINUTES AGO TIMESERIES
```

0 = Port closed
1 = Port open

## Building

Golang is required to build the integration. We recommend Golang 1.11 or higher.

After cloning this repository, go to the directory of the Port Monitor integration and build it:

```bash
$ make
```

The command above executes the tests for the Port Monitor integration and builds an executable file called `nri-port-monitor` under the `bin` directory.  

To start the integration, run `nri-port-monitor`:

```bash
$ ./bin/nri-port-monitor
```

If you want to know more about usage of `./bin/nri-port-monitor`, pass the `-help` parameter:

```bash
$ ./bin/nri-port-monitor -help
```

External dependencies are managed through the [govendor tool](https://github.com/kardianos/govendor). Locking all external dependencies to a specific version (if possible) into the vendor directory is required.

## Testing

To run the tests execute:

```bash
$ make test
```

## Support

Should you need assistance with New Relic products, you are in good hands with several support diagnostic tools and support channels.

> This [troubleshooting framework](https://discuss.newrelic.com/t/troubleshooting-frameworks/108787) steps you through common troubleshooting questions.

> New Relic offers NRDiag, [a client-side diagnostic utility](https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/troubleshooting/new-relic-diagnostics) that automatically detects common problems with New Relic agents. If NRDiag detects a problem, it suggests troubleshooting steps. NRDiag can also automatically attach troubleshooting data to a New Relic Support ticket.

If the issue has been confirmed as a bug or is a Feature request, please file a Github issue.

**Support Channels**

* [New Relic Documentation](https://docs.newrelic.com): Comprehensive guidance for using our platform
* [New Relic Community](https://discuss.newrelic.com): The best place to engage in troubleshooting questions
* [New Relic Developer](https://developer.newrelic.com/): Resources for building a custom observability applications
* [New Relic University](https://learn.newrelic.com/): A range of online training for New Relic users of every level

## Privacy

At New Relic we take your privacy and the security of your information seriously, and are committed to protecting your information. We must emphasize the importance of not sharing personal data in public forums, and ask all users to scrub logs and diagnostic information for sensitive information, whether personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or identifiable individual, including, for example, your name, phone number, post code or zip code, Device ID, IP address and email address.

Review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy) for more information.

## Contributing

We encourage your contributions to improve the Port Monitor integration! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.

If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](/SECURITY.md), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, please review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project would not be what it is today.

## License

nri-port-monitor is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
