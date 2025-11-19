# --- Builder Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod ./
RUN go get -d ./...
RUN go mod tidy
RUN go mod download
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get gopkg.in/yaml.v3

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o awair-exporter .

# --- Final Stage ---
FROM alpine:latest

# Install necessary packages
# ca-certificates: for making HTTPS calls (if ever needed)
# su-exec: for dropping root privileges
RUN apk --no-cache add ca-certificates su-exec

# Create a directory for configuration
RUN mkdir /config

# Copy the application binary from the builder stage
COPY --from=builder /app/awair-exporter /usr/local/bin/awair-exporter

# Copy and set up the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose the default port for Prometheus scrapes
EXPOSE 9101

# Set the entrypoint
ENTRYPOINT ["/entrypoint.sh"]

# Set the default command to run the exporter
CMD ["/usr/local/bin/awair-exporter"]
