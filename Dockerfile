# Build Stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy dependency files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/adzan app/main.go

# Final Runtime Stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/adzan /app/adzan

# Directory for external data mount
RUN mkdir -p /data
VOLUME ["/data"]

# Default environment configuration for Docker
ENV DATA_SOURCE=csv
ENV CSV_PATH=/data
ENV SERVER_PORT=4321
ENV SERVER_LOG_LEVEL=info

# Expose HTTP and gRPC ports
EXPOSE 4321
EXPOSE 54321

ENTRYPOINT ["/app/adzan"]
