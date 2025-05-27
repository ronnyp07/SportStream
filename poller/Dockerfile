# Build stage
FROM golang:1.23-alpine AS builder

RUN go env -w GOPRIVATE='github.com/sts-solutions'

# Copy only poller folder (since it has its own go.mod)
COPY . /app
WORKDIR /app

# Build the binary
RUN go build -o poller_bin ./cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /opt/poller

# Copy binary and config
COPY --from=builder /app/poller_bin ./bin/main
COPY --from=builder /app/config/app/config.yaml ./config/app/config.yaml
COPY --from=builder /app/infra.env ./infra.env

ENTRYPOINT ["./bin/main", "-conf", "/opt/poller/config/app/config.yaml"]

EXPOSE 80
