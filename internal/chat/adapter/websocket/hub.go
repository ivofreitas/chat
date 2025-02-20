package websocket

import (
	"github.com/ivofreitas/chat/internal/chat/adapter/broker"
	"github.com/ivofreitas/chat/pkg/log"
	"sync"
)

var (
	once     sync.Once
	instance *Hub
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	Broadcast chan []byte
	Register  chan *Client

	clients    map[*Client]bool
	unregister chan *Client
	publisher  broker.Publisher
}

func GetHub() *Hub {
	once.Do(func() {
		instance = &Hub{
			Broadcast:  make(chan []byte),
			Register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			publisher:  broker.NewPublisher(),
		}
	})
	return instance
}

func (h *Hub) Run() {
	for {
		entry := log.NewEntry()
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			go h.publisher.Publish(message)
			entry.Infof("received a message!: %s", string(message))
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
