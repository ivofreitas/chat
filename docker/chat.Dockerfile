FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o chat ./cmd/chat

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/chat .
CMD ["./chat"]
