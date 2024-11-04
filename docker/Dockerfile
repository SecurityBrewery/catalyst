FROM ubuntu:24.04

RUN apt-get update && apt-get install -y curl python3 python3-pip python3-venv

COPY catalyst /usr/local/bin/catalyst

EXPOSE 8080

VOLUME /usr/local/bin/catalyst_data

HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/usr/local/bin/catalyst", "serve", "--http", "0.0.0.0:8080"]