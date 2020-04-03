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
  -default-uptime-slo float
    	default uptime SLO to be used when the check doesn't provide a uptime SLO tag (i.e. uptime_slo_999 to 99.9% uptime SLO) (default 99)
  -outage-check-period int
    	time (in days) in which to retrieve outage data from the Pingdom API (default 7)
  -port int
    	port to listen on (default 9158)
```

#### Supported Pingdom Tags

##### `uptime_slo_xxx`

This will instruct pingdom-exporter to use a custom SLO for the given check
instead of the default one of 99%. Some tag examples and their corresponding
SLOs:

- `uptime_slo_99` - 99%, same as default
- `uptime_slo_995` - 99.5%
- `uptime_slo_999` - 99.9%

##### `pingdom_exporter_ignored`

Checks with this tag won't have their metrics exported. Use this when you don't
want to disable some check just to have it excluded from the pingdom-exporter
metrics.

(TODO: We might change the implementation for this later; one idea is to allow
the pingdom-exporter to only operate in checks that contain a specific tag -
that way, we can have different exporters targeting different checks based on
their label)

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

| Metric Name                                         | Description                                                                     |
| --------------------------------------------------- | ------------------------------------------------------------------------------- |
| `pingdom_up`                                        | Was the last query on Pingdom API successful                                    |
| `pingdom_uptime_status`                             | The current status of the check (1: up, 0: down)                                |
| `pingdom_uptime_response_time_seconds`              | The response time of last test, in seconds                                      |
| `pingdom_slo_period_seconds`                        | Outage check period, in seconds (see `-outage-check-period` flag)               |
| `pingdom_outages_total`                             | Number of outages within the outage check period                                |
| `pingdom_down_seconds`                              | Total down time within the outage check period, in seconds                      |
| `pingdom_up_seconds`                                | Total up time within the outage check period, in seconds                        |
| `pingdom_uptime_slo_error_budget_total_seconds`     | Maximum number of allowed downtime, in seconds, according to the uptime SLO     |
| `pingdom_uptime_slo_error_budget_available_seconds` | Number of seconds of downtime we can still have without breaking the uptime SLO |

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
