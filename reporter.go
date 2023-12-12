package goerrorcollector

import (
	"context"
	"encoding/json"
	"github.com/latifrons/latigo/mq/publisher"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"sync"
)

type Option func(*GoErrorCollector)

type ErrorMessage struct {
	Component  string `json:"component"`
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
	Time       int64  `json:"time"`
	Severity   int    `json:"severity"`
}

type ReceiverMqOption struct {
	RabbitMQUrl  string
	ExchangeName string
	Topic        string
}

type GoErrorCollector struct {
	start            sync.Once
	receiverMqOption ReceiverMqOption
	publisher        *publisher.ReliableRabbitPublisher
	buffer           chan ErrorMessage
	quitChan         chan bool
}

func WithReceiverRabbitMQ(option ReceiverMqOption) Option {
	return func(collector *GoErrorCollector) {
		collector.receiverMqOption = option
	}
}

func (g *GoErrorCollector) Start(options ...Option) (err error) {
	g.start.Do(func() {
		g.buffer = make(chan ErrorMessage, 100)
		g.quitChan = make(chan bool)

		for _, option := range options {
			option(g)
		}
		err = g.startReporting()
		if err != nil {
			log.Error().Err(err).Msg("failed to start reporting")
		}
	})
	return
}

func (g *GoErrorCollector) Stop() {
	g.quitChan <- true
}

func (g *GoErrorCollector) Name() string {
	return "GoErrorCollector"
}

func (g *GoErrorCollector) initPublisher(channel *amqp091.Channel) (err error) {
	err = channel.ExchangeDeclare(g.receiverMqOption.ExchangeName, "topic", true, false, false, false, nil)
	return
}

func (g *GoErrorCollector) startReporting() (err error) {
	g.publisher = publisher.NewReliableRabbitPublisher(g.receiverMqOption.RabbitMQUrl, publisher.WithInitFunc(g.initPublisher))
	err = g.publisher.Start()
	if err != nil {
		return
	}
	go g.sending()
	return
}

func (g *GoErrorCollector) Report(msg ErrorMessage) {
	select {
	case g.buffer <- msg:
		return
	default:
		log.Debug().Msg("buffer is full")
	}
}

func (g *GoErrorCollector) sending() {
	for {
		select {
		case msg := <-g.buffer:
			j, err := json.Marshal(msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal message in error collector")
				continue
			}
			err = g.publisher.Publish(context.Background(), g.receiverMqOption.ExchangeName, g.receiverMqOption.Topic, amqp091.Publishing{
				Body: j,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to publish message in error collector")
				continue
			}
		case _ = <-g.quitChan:
			return
		}
	}
}
