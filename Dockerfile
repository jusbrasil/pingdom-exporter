FROM centurylink/ca-certs
MAINTAINER Daniel Martins <daniel.martins@jusbrasil.com.br>

COPY ./bin/pingdom-exporter /pingdom-exporter
ENTRYPOINT ["/pingdom-exporter"]

USER 65534:65534
