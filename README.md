# Pingdom exporter

[![Docker Pulls](https://img.shields.io/docker/pulls/vptech/pingdom-exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/strike-team/pingdom_exporter)][goreportcard]

Prometheus exporter for uptime and transaction metrics exposed by Pingdom API, written in Go.

To run it:

```bash
make
./pingdom_exporter server <pingdom_username> <pingdom_password> <pingdom_token>
```

## Exported Metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| pingdom_up | Was the last query on Pingdom API successful, | |
| pingdom_uptime_status | The current status of the check (1: up, 0: down). | |
| pingdom_uptime_response_time | The response time of last test in milliseconds. | |
| pingdom_transaction_status | The current status of the transaction (1: successful, 0: failing). | |

## Using Docker

You can deploy this exporter using the [vptech/pingdom-exporter](https://hub.docker.com/r/vptech/pingdom-exporter/) Docker image.

For example:

```bash
docker pull vptech/pingdom-exporter

docker run -d -p 9158:9158 \
        vptech/pingdom-exporter \
        <pingdom_username> \
        <pingdom_password> \
        <pingdom_token>
```

[hub]: https://hub.docker.com/r/vptech/pingdom-exporter/
[goreportcard]: https://goreportcard.com/report/github.com/strike-team/pingdom_exporter