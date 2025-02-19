package application

import (
	"context"
	"github.com/ivofreitas/chat/internal/bot/adapter/broker"
	client "github.com/ivofreitas/chat/internal/bot/adapter/external_api"
	"github.com/ivofreitas/chat/pkg/config"
	"github.com/ivofreitas/chat/pkg/log"
	"github.com/ivofreitas/chat/pkg/rabbitmq"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

type Receiver struct {
	stockAPI  client.StockAPI
	publisher *broker.Publisher
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
	logger    *logrus.Entry
	signal    chan struct{}
}

var stockPattern = regexp.MustCompile(`/stock=[a-zA-Z0-9_]+`)

func NewReceiver() *Receiver {
	env := config.GetEnv()

	conn, ch := rabbitmq.Get()

	return &Receiver{
		stockAPI:  client.NewClient(env.External.StockBaseUrl),
		publisher: broker.NewPublisher(),
		conn:      conn, ch: ch,
		queueName: env.Broker.StockCodeQueue,
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
		entry := log.NewEntry()

		for msg := range messages {
			pattern := string(msg.Body)
			if !stockPattern.MatchString(pattern) {
				continue
			}

			var code string
			if split := strings.Split(pattern, "="); len(split) == 2 {
				code = split[1]
			} else {
				continue
			}

			stocks, err := r.stockAPI.Lookup(context.Background(), code)
			if err != nil {
				entry.Errorf("lookup error: %s", err.Error())
			}

			for _, stock := range stocks {
				r.publisher.Publish([]byte(stock.String()))
			}
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
