FROM golang:1.24.1 AS builder
LABEL authors="rtukpe"

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o compareids ./cmd

FROM alpine:latest

COPY --from=builder /go/src/app/compareids /usr/local/bin/compareids
COPY --from=builder /go/src/app/run.sh /run.sh

RUN chmod +x /run.sh

# Install necessary dependencies
RUN apk add --no-cache git make

CMD ["compareids"]