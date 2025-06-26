# Start from the official Golang image for building
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o game-integration-api ./cmd

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/game-integration-api .
EXPOSE 8080
CMD ["./game-integration-api"] 