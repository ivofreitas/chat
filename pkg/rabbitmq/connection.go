package rabbitmq

import (
	"github.com/ivofreitas/chat/pkg/config"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	once sync.Once
)

// Get establishes a RabbitMQ connection and channel
func Get() (*amqp.Connection, *amqp.Channel) {
	once.Do(func() {
		var err error

		rabbitURL := config.GetEnv().Broker.Url

		conn, err = amqp.Dial(rabbitURL)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}

		ch, err = conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}

		// Ensure the queues exist
		queues := []string{
			config.GetEnv().Broker.StockCodeQueue,
			config.GetEnv().Broker.StockQuoteQueue,
		}

		for _, queue := range queues {
			_, err := ch.QueueDeclare(
				queue, // name
				true,  // durable
				false, // auto-delete
				false, // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				log.Fatalf("Failed to declare queue %s: %v", queue, err)
			}
		}
	})

	return conn, ch
}
