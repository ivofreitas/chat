package application

import (
	"github.com/ivofreitas/chat/internal/chat/adapter/websocket"
	"github.com/ivofreitas/chat/pkg/config"
	"github.com/ivofreitas/chat/pkg/log"
	"github.com/ivofreitas/chat/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Receiver struct {
	hub       *websocket.Hub
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
	logger    *logrus.Entry
	signal    chan struct{}
}

func NewReceiver() *Receiver {
	env := config.GetEnv()

	conn, ch := rabbitmq.Get()

	return &Receiver{
		hub:  websocket.GetHub(),
		conn: conn, ch: ch,
		queueName: env.Broker.StockQuoteQueue,
		logger:    log.NewEntry(),
		signal:    make(chan struct{})}
}

func (r *Receiver) Run() {
	r.consume()
	r.logger.Println("Receiver started and waiting for the graceful signal...")
	r.watchStop()
}

func (r *Receiver) consume() {

	messages, err := r.ch.Consume(
		r.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.logger.WithError(err).Fatal("Shutting down the receiver now")
	}

	go func() {
		for msg := range messages {
			r.hub.Broadcast <- msg.Body
		}
	}()
}

func (s *Receiver) watchStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	s.ch.Close()
	s.conn.Close()

	s.logger.Info("Receiver is stopping...")

	close(s.signal)
}
