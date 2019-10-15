# Pingdom Metrics Exporter for Prometheus

Prometheus exporter for uptime and transaction metrics exposed by Pingdom API.

To run it:

```sh
make

# Provide the API token via an environment variable
export PINGDOM_API_TOKEN=<api-token>

bin/pingdom-exporter
```

## Options

```sh
$ bin/pingdom-exporter -h
Usage of bin/pingdom-exporter:
  -outage-check-period int
    	time (in days) in which to retrieve outage data from the Pingdom API (default 7)
  -port int
    	port to listen on (default 9158)
  -wait int
    	time (in seconds) between accessing the Pingdom  API (default 60)
```

## Exported Metrics

| Metric                                 | Description                                                          |
| -------------------------------------- | -------------------------------------------------------------------- |
| `pingdom_up`                           | Was the last query on Pingdom API successful.                        |
| `pingdom_uptime_status`                | The current status of the check (1: up, 0: down).                    |
| `pingdom_uptime_response_time_seconds` | The response time of last test, in seconds.                          |
| `pingdom_outage_check_period_seconds`  | Outage check period, in seconds (see the -outage-check-period flag). |
| `pingdom_outages_total`                | Number of outages within the outage check period.                    |
| `pingdom_down_seconds`                 | Total down time within the outage check period.                      |
| `pingdom_up_seconds`                   | Total up time within the outage check period.                        |

## Using Docker

You can deploy this exporter using the
[jusbrasil/pingdom-exporter](https://hub.docker.com/r/jusbrasil/pingdom-exporter/)
Docker image:

```bash
docker run -d -p 9158:9158 \
        -e PINGDOM_API_TOKEN=<api-token> \
        jusbrasil/pingdom-exporter
```
