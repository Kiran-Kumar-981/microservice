# Stage 1: Build the Go app
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o server ./cmd/web

# Stage 2: Create a minimal image for the Go app
FROM alpine:latest

WORKDIR /server/

COPY --from=builder /app/server .

CMD ["./server"]
