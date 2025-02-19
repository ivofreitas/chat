# Chat Application

[![Go Version](https://img.shields.io/badge/Go-1.23.5-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/ivofreitas/chat/actions)
[![License](https://img.shields.io/badge/license-Apache_2.0-blue.svg)](LICENSE)

## Overview
Chat Application is a real-time messaging platform built with Golang, supporting user authentication, WebSocket-based messaging, and scalable architecture.

## Table of Contents
- [Makefile Commands](#makefile-commands)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [WebSocket Integration](#websocket-integration)
- [API Documentation](#api-documentation)
- [Future Improvements](#future-improvements)
- [License](#license)

## Makefile Commands
The project includes a `Makefile` that simplifies common development tasks. The key commands are:

| Command            | Description                                         |
|--------------------|-----------------------------------------------------|
| `make all`         | Run all tests, then build and run                   |
| `make build`       | Build all binaries (`chat-server`, `auth-server`, `bot-server`) |
| `make clean`       | Remove build artifacts and tidy up dependencies     |
| `make run`         | Build and run all servers (`chat`, `auth`, `bot`)  |
| `make docker-up`   | Start the application with Docker Compose           |
| `make docker-down` | Stop and remove the Docker Compose containers       |
| `make lint`        | Run code linters                                    |
| `make test`        | Run all unit tests with race detection and coverage |
| `make mock`        | Generate mocks using Mockery                        |
| `make swag`        | Generate API documentation using Swagger            |
| `make help`        | Display available commands                          |

## Environment Variables
The following environment variables are used in the application:

| Name                        | Suggested Value                            | Required |
|-----------------------------|--------------------------------------------|----------|
| `SERVER_CHAT_PORT`          | `8081`                                     | ✅       |
| `SERVER_AUTH_PORT`          | `8082`                                     | ✅       |
| `LOG_ENABLED`               | `true`                                     | ✅       |
| `LOG_LEVEL`                 | `info`                                     | ✅       |
| `DB_HOST`                   | `localhost`                                | ✅       |
| `DB_PORT`                   | `5432`                                     | ✅       |
| `DB_USER`                   | `chat_user`                                | ✅       |
| `DB_PASSWORD`               | `chat_password`                            | ✅       |
| `DB_NAME`                   | `chat_db`                                  | ✅       |
| `DB_SSLMODE`                | `disable`                                  | ❌       |
| `BROKER_URL`                | `amqp://user:password@localhost:5672/`     | ✅       |
| `BROKER_USER`               | `user`                                     | ✅       |
| `BROKER_PASSWORD`           | `password`                                 | ✅       |
| `BROKER_STOCK_CODE_QUEUE`   | `stock_code_queue`                         | ✅       |
| `BROKER_STOCK_QUOTE_QUEUE`  | `stock_quote_queue`                        | ✅       |
| `EXTERNAL_STOCK_BASE_URL`   | `https://stooq.com/q/l/`                   | ✅       |
| `SECURITY_JWT_SECRET_KEY`   | `701b86a867843ba2ab7098140bdff88de0ee0970f09e60ef4a22196eb36237dc` | ✅       |

## API Endpoints

### Users Group
| Method   | Endpoint          | Description                   |
|----------|------------------|-------------------------------|
| `POST`   | `/users/register` | Register a new user          |
| `POST`   | `/users/login`    | Authenticate user            |
| `GET`    | `/users/profile`  | Get user profile             |

## WebSocket Integration
The Chat Application uses WebSockets for real-time messaging. Connect to:
```
ws://localhost:8081/ws
```

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
