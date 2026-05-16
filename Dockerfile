# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/customer-support-hub ./cmd/api

FROM alpine:3.21
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /bin/customer-support-hub /usr/local/bin/customer-support-hub
COPY .env.example ./.env.example
USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/customer-support-hub"]
