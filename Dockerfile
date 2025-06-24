# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY internal/ internal
RUN go mod download

COPY main/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/rssagg .

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rssagg .

ENTRYPOINT ["/app/rssagg"]