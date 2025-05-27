# Build stage
FROM golang:1.23-alpine AS builder

RUN go env -w GOPRIVATE='github.com/sts-solutions'

# Copy only api folder (since it has its own go.mod)
COPY . /app
WORKDIR /app

# Build the binary
RUN go build -o api_bin ./cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /opt/api

# Copy binary and config
COPY --from=builder /app/api_bin ./bin/main
COPY --from=builder /app/config/app/config.yaml ./config/app/config.yaml
COPY --from=builder /app/infra.env ./infra.env

ENTRYPOINT ["./bin/main", "-conf", "/opt/api/config/app/config.yaml"]

EXPOSE 80
