# Pingdom Metrics Exporter for Prometheus

Prometheus exporter for uptime and transaction metrics exposed by Pingdom API.

To run it:

```bash
make

# Provide the API token via an environment variable
export PINGDOM_API_TOKEN=<api-token>

./pingdom_exporter server
```

## Exported Metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| pingdom_up | Was the last query on Pingdom API successful, | |
| pingdom_uptime_status | The current status of the check (1: up, 0: down). | |
| pingdom_uptime_response_time | The response time of last test in milliseconds. | |

## Using Docker

You can deploy this exporter using the
[jusbrasil/pingdom-exporter](https://hub.docker.com/r/jusbrasil/pingdom-exporter/)
Docker image:

```bash
docker run -d -p 9158:9158 \
        -e PINGDOM_API_TOKEN=<api-token> \
        jusbrasil/pingdom-exporter \
        <pingdom_username> \
        <pingdom_password> \
        <pingdom_token>
```
