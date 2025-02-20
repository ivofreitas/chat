package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/ivofreitas/chat/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(message []byte) {
	m.Called(message)
}

func TestWebsocket(t *testing.T) {
	mockPublisher := new(MockPublisher)
	hub := &Hub{
		Broadcast:  make(chan []byte, 1),
		Register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		clients:    make(map[*Client]bool),
		publisher:  mockPublisher,
	}

	log.Init()

	go hub.Run()
	time.Sleep(200 * time.Millisecond)

	testCases := []struct {
		name     string
		setup    func(hub *Hub, client *Client)
		validate func(t *testing.T, hub *Hub, client *Client)
	}{
		{
			name: "Add Client",
			setup: func(hub *Hub, client *Client) {
				hub.Register <- client
				time.Sleep(100 * time.Millisecond)
			},
			validate: func(t *testing.T, hub *Hub, client *Client) {
				assert.Contains(t, hub.clients, client)
			},
		},
		{
			name: "Remove Client",
			setup: func(hub *Hub, client *Client) {
				hub.Register <- client
				time.Sleep(100 * time.Millisecond)
				hub.unregister <- client
				time.Sleep(100 * time.Millisecond)
			},
			validate: func(t *testing.T, hub *Hub, client *Client) {
				assert.NotContains(t, hub.clients, client)
			},
		},
		{
			name: "Broadcast Message",
			setup: func(hub *Hub, client *Client) {
				hub.Register <- client
				time.Sleep(100 * time.Millisecond)
				mockPublisher.On("Publish", mock.Anything).Return(nil)
				hub.Broadcast <- []byte("Broadcast Message")
				time.Sleep(200 * time.Millisecond)
			},
			validate: func(t *testing.T, hub *Hub, client *Client) {
				mockPublisher.AssertCalled(t, "Publish", []byte("Broadcast Message"))
				mockPublisher.AssertExpectations(t)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &Client{Hub: hub, Conn: &websocket.Conn{}, Send: make(chan []byte, 1)}
			if tc.setup != nil {
				tc.setup(hub, client)
			}
			if tc.validate != nil {
				tc.validate(t, hub, client)
			}
		})
	}
}
