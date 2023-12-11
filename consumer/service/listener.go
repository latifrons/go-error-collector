package service

import (
	"context"
	"github.com/latifrons/latigo/mq/consumer"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type ErrorMessageHandler interface {
	HandleErrorMessage(amqp091.Delivery) (err error)
}

type ErrorListener struct {
	ErrorMessageHandler ErrorMessageHandler `container:"type"`
	RabbitMQUrl         string
	ExchangeName        string
	QueuePrefetchCount  int
	queueName           string
	routingKey          string

	consumer *consumer.ReliableRabbitConsumer
}

func (e *ErrorListener) InitDefault() {
	e.queueName = "error_delivery"
	e.routingKey = "*"
}

func (c *ErrorListener) Start() {
	c.consumer = consumer.NewReliableRabbitConsumer(c.RabbitMQUrl,
		c.ExchangeName, c.queueName, c.routingKey,
		c.handler,
		consumer.WithInitFunc(c.initConsumer),
		consumer.WithCleanFunc(c.cleanConsumer),
		consumer.WithQos(c.QueuePrefetchCount, false))
	err := c.consumer.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start consumer")
	}
}

func (e *ErrorListener) Stop() {
	//TODO implement me
	panic("implement me")
}

func (e *ErrorListener) Name() string {
	return "ErrorListener"
}

func (c *ErrorListener) initConsumer(channel *amqp091.Channel) (err error) {
	declare, err := channel.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = channel.QueueBind(declare.Name, c.routingKey, c.ExchangeName, false, nil)
	if err != nil {
		return err
	}
	return
}

func (c *ErrorListener) cleanConsumer(channel *amqp091.Channel) (err error) {
	_, err = channel.QueueDelete(c.queueName, false, false, false)
	return
}

func (c *ErrorListener) handler(ctx context.Context, msgBody amqp091.Delivery) interface{} {
	err := c.ErrorMessageHandler.HandleErrorMessage(msgBody)
	if err != nil {
		log.Error().Err(err).Uint64("i", msgBody.DeliveryTag).Msg("failed to handle message")
		//err = msgBody.Nack(false, true)
		//if err != nil {
		//	log.Error().Err(err).Uint64("i", msgBody.DeliveryTag).Msg("failed to nack message")
		//}
	} else {
		err = msgBody.Ack(false)
		if err != nil {
			log.Error().Err(err).Uint64("i", msgBody.DeliveryTag).Msg("failed to ack message")
		}
	}
	return nil
}
