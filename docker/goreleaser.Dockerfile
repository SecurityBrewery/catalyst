FROM ubuntu:24.04

COPY catalyst /usr/local/bin/catalyst

EXPOSE 8080

VOLUME /usr/local/bin/catalyst_data

CMD ["/usr/local/bin/catalyst", "serve", "--http", "0.0.0.0:8080"]