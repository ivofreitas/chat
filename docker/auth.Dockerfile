FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o auth ./cmd/auth

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth .
CMD ["./auth"]
