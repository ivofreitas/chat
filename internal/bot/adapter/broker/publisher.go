package broker

import (
	"github.com/ivofreitas/chat/pkg/config"
	"github.com/ivofreitas/chat/pkg/log"
	"github.com/ivofreitas/chat/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Publisher struct {
	ch        *amqp.Channel
	queueName string
	logger    *logrus.Entry
}

func NewPublisher() *Publisher {
	env := config.GetEnv()

	_, ch := rabbitmq.Get()

	return &Publisher{ch: ch, queueName: env.Broker.StockQuoteQueue, logger: log.NewEntry()}
}

func (p *Publisher) Publish(message []byte) {
	q, err := p.ch.QueueDeclare(
		p.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		p.logger.Error(err)
	}

	err = p.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		p.logger.Error(err)
	}
}
