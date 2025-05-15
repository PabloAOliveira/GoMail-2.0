FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o email-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/email-service .
COPY --from=builder /app/.env .env

EXPOSE 2000
CMD ["./email-service"]

