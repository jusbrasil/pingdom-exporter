# Pingdom Metrics Exporter for Prometheus

Prometheus exporter for uptime metrics exposed by the Pingdom API.

## Running

Make sure you expose the Pingdom API Token via the `PINGDOM_API_TOKEN`
environment variable:

```sh
# Expose the Pingdom API Token
export PINGDOM_API_TOKEN=<api-token>

# Run the binary with the default options
bin/pingdom-exporter
```

### Usage

```
bin/pingdom-exporter -h

Usage of bin/pingdom-exporter:
  -outage-check-period int
    	time (in days) in which to retrieve outage data from the Pingdom API (default 7)
  -port int
    	port to listen on (default 9158)
  -wait int
    	time (in seconds) to wait between each metrics update (default 60)
```

### Docker Image

You can run this exporter using the
[jusbrasil/pingdom-exporter](https://hub.docker.com/r/jusbrasil/pingdom-exporter/)
Docker image:

```bash
docker run -d -p 9158:9158 \
        -e PINGDOM_API_TOKEN=<api-token> \
        jusbrasil/pingdom-exporter
```

## Exported Metrics

| Metric Name                            | Description                                                       |
| -------------------------------------- | ----------------------------------------------------------------- |
| `pingdom_up`                           | Was the last query on Pingdom API successful                      |
| `pingdom_uptime_status`                | The current status of the check (1: up, 0: down)                  |
| `pingdom_uptime_response_time_seconds` | The response time of last test, in seconds                        |
| `pingdom_outage_check_period_seconds`  | Outage check period, in seconds (see `-outage-check-period` flag) |
| `pingdom_outages_total`                | Number of outages within the outage check period                  |
| `pingdom_down_seconds`                 | Total down time within the outage check period, in seconds        |
| `pingdom_up_seconds`                   | Total up time within the outage check period, in seconds          |

## Development

All relevant commands are exposed via Makefile targets:

```sh
# Build the binary
make

# Run the tests
make test

# Check linting rules
make lint

# Build Docker image
make docker-build

# Push Docker images to registry
make docker-push
```
