FROM centurylink/ca-certs
MAINTAINER Daniel Martins <daniel.martins@jusbrasil.com.br>

COPY ./bin/pingdom_exporter /pindom_exporter
ENTRYPOINT ["/pingdom_exporter"]

USER 65534:65534
