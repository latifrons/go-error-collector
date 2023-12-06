package report

import (
	go_error_collector "github.com/latifrons/go-error-collector"
	"time"
)

var Reporter = go_error_collector.GoErrorCollector{}

func Report(component string, severity int, message string, stacktrace string) {
	Reporter.Report(go_error_collector.Message{
		Component:  component,
		Message:    message,
		Stacktrace: stacktrace,
		Time:       time.Now().UnixMilli(),
		Severity:   severity,
	})
}

func Start(url string, exchange string, topic string) {
	Reporter.Start(go_error_collector.WithReceiverRabbitMQ(go_error_collector.RecieverMqOption{
		RabbitMQUrl:  url,
		ExchangeName: exchange,
		Topic:        topic,
	}))
}
