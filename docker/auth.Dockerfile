FROM golang:1.23.5-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download
RUN go build -o auth_service cmd/auth/main.go

FROM alpine as release

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=builder /app/auth_service /app/auth_service
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

EXPOSE 8081

ENTRYPOINT ["/app/wait-for-it.sh", "postgres", "5432", "--", "./auth_service"]
