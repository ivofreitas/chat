package application

import (
	"github.com/ivofreitas/chat/internal/auth/application/user"
	"github.com/ivofreitas/chat/pkg/postgres"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

func register(echo *echo.Echo) {
	userGroup(echo)
	swaggerGroup(echo)
}

func swaggerGroup(echo *echo.Echo) {
	echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

func userGroup(e *echo.Echo) {

	repository := user.NewRepository(postgres.NewConnection())
	handler := user.NewHandler(repository)

	users := e.Group("/users")
	users.POST("/register", handler.RegisterUser)
	users.POST("/login", handler.LoginUser)
}
