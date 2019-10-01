FROM alpine:3.8
MAINTAINER Joseph Salisbury <joseph@giantswarm.io>

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/pingdom_exporter   /bin/pingdom_exporter

RUN apk update && apk add ca-certificates
EXPOSE     9652
ENTRYPOINT [ "/bin/pingdom_exporter" ]