FROM golang:1.24.1 AS builder
LABEL authors="rtukpe"

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o compareids ./cmd

FROM alpine:latest

# Install necessary dependencies for large-scale testing
RUN apk add --no-cache \
    git \
    make \
    postgresql-client \
    postgresql-dev \
    build-base \
    curl \
    jq \
    bash \
    coreutils \
    procps \
    && rm -rf /var/cache/apk/*

# Create directories for results and configuration
RUN mkdir -p /app/results /app/config

COPY --from=builder /go/src/app/compareids /usr/local/bin/compareids
COPY --from=builder /go/src/app/run.sh /run.sh
COPY --from=builder /go/src/app/config/ /app/config/

RUN chmod +x /run.sh

# Set working directory
WORKDIR /app

# Default environment variables for configuration
ENV TEST_SCENARIO=medium
ENV MAX_ROWS=1000000
ENV PARALLEL_TESTS=1
ENV RESULTS_DIR=/app/results
ENV CONFIG_DIR=/app/config

CMD ["/run.sh"]