ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/pingdom_exporter   /bin/pingdom_exporter

EXPOSE     9652
ENTRYPOINT [ "/bin/pingdom_exporter" ]