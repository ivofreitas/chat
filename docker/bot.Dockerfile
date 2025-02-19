FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o bot ./cmd/bot

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bot .
CMD ["./bot"]
