FROM golang:1.22

WORKDIR /app
ADD . .
RUN make build

ENTRYPOINT ["/app/bin/pingdom-exporter"]

USER 65534:65534
