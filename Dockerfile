FROM golang:1.22 AS build

WORKDIR /app
ADD . .
RUN make build

FROM alpine:latest
MAINTAINER Daniel Martins <daniel.martins@jusbrasil.com.br>

RUN apk add --no-cache ca-certificates \
  && update-ca-certificates
COPY --from=build /app/bin/pingdom-exporter /pingdom-exporter
ENTRYPOINT ["/pingdom-exporter"]

USER 65534:65534
