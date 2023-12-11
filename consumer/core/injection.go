package core

import (
	"github.com/golobby/container/v3"
	"github.com/latifrons/commongo/safe_viper"
	"github.com/latifrons/goerrorcollector/consumer/service"
)

type SingletonBatch []interface{}

func BuildDependencies() (err error) {
	singletons := []interface{}{
		func() *service.ErrorListener {
			var c service.ErrorListener
			err := container.Fill(&c)
			if err != nil {
				panic(err)
			}
			c.RabbitMQUrl = safe_viper.ViperMustGetString("listener.rabbitmq_url")
			c.ExchangeName = safe_viper.ViperMustGetString("listener.exchange_name")
			c.QueuePrefetchCount = 100
			return &c
		},
		func() *service.ErrorMessageHandler {
			var c service.ErrorMessageHandler
			err := container.Fill(&c)
			if err != nil {
				panic(err)
			}
			return &c
		},
	}
	for _, singleton := range singletons {
		err = container.SingletonLazy(singleton)
		if err != nil {
			return
		}

	}
	return
}
