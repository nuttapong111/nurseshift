## Multi-stage Dockerfile for user-service

FROM golang:1.22-alpine AS build
WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=build /app/server ./server

EXPOSE 8082

CMD ["./server"]


