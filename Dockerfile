FROM golang:1.13 AS build

WORKDIR /app
ADD . .
RUN make test lint build

FROM centurylink/ca-certs
MAINTAINER Daniel Martins <daniel.martins@jusbrasil.com.br>

COPY --from=build /app/bin/pingdom-exporter /pingdom-exporter
ENTRYPOINT ["/pingdom-exporter"]

USER 65534:65534
