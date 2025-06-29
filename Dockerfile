FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o game-integration-api ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/game-integration-api .
EXPOSE 8080
CMD ["./game-integration-api"]

FROM golang:1.24-alpine AS dev
WORKDIR /app
ENV GOPROXY=direct
RUN go install github.com/cosmtrek/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["air"] 