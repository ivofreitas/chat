package application

import (
	gorillaws "github.com/gorilla/websocket"
	"github.com/ivofreitas/chat/internal/chat/adapter/websocket"
	"github.com/ivofreitas/chat/pkg/log"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	upgrader = gorillaws.Upgrader{}
)

func register(echo *echo.Echo) {
	wsGroup(echo)
}

func wsGroup(e *echo.Echo) {

	e.GET("/ws", func(c echo.Context) error {
		entry := log.NewEntry()
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			entry.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		client := &websocket.Client{Hub: websocket.GetHub(), Conn: conn, Send: make(chan []byte, 256)}
		client.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump()

		return nil
	})

	e.Static("/", "static")

	e.GET("/", func(c echo.Context) error {
		if c.Request().URL.Path != "/" {
			return echo.NewHTTPError(http.StatusNotFound, "Not found")
		}

		if c.Request().Method != http.MethodGet {
			return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
		}

		return c.File("static/login.html")
	})
}
