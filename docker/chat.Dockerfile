FROM golang:1.23.5-alpine as builder

WORKDIR /app

COPY . .

RUN go get github.com/swaggo/swag/gen@v1.16.4
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN swag init -g cmd/chat/main.go
RUN go mod tidy
RUN go mod download
RUN go build -o chat_service cmd/chat/main.go

FROM alpine as release

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=builder /app/chat_service /app/chat_service

COPY static /app/static

COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

EXPOSE 8082

ENTRYPOINT ["/app/wait-for-it.sh", "rabbitmq", "5672", "--", "./chat_service"]
