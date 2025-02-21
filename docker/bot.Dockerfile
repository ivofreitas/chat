FROM golang:1.23.5-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download
RUN go build -o bot_service cmd/bot/main.go

FROM alpine as release

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=builder /app/bot_service /app/bot_service
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

ENTRYPOINT ["/app/wait-for-it.sh", "rabbitmq", "5672", "--", "./bot_service"]
