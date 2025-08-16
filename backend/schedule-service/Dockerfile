## Multi-stage Dockerfile for schedule-service

# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /app

# Enable Go modules and caching
ENV CGO_ENABLED=0

# Pre-copy mod files to leverage Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app

# Add timezone data and certs
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from build stage
COPY --from=build /app/server ./server

# Railway will inject PORT at runtime; expose a default for local usage
EXPOSE 8084

# Start the service
CMD ["./server"]


