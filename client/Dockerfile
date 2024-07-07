# Stage 1: Build the Go app
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client ./cmd/api

# Stage 2: Create a minimal image for the Go app
FROM alpine:latest

WORKDIR /client/

COPY --from=builder /app/client .

CMD ["./client"]
