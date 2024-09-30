FROM oven/bun:debian
RUN apt-get update && apt-get install -y make
COPY .. /tmp/catalyst

WORKDIR /tmp/catalyst

RUN make build-ui

FROM golang:1.23
COPY --from=0 /tmp/catalyst /tmp/catalyst

WORKDIR /tmp/catalyst

RUN go build -o /usr/local/bin/catalyst

FROM ubuntu:24.04

COPY --from=1 /usr/local/bin/catalyst /usr/local/bin/catalyst

EXPOSE 8080

VOLUME /usr/local/bin/catalyst_data

CMD ["/usr/local/bin/catalyst", "serve", "--http", "0.0.0.0:8080"]