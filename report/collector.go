package report

import (
	"github.com/latifrons/goerrorcollector"
	"time"
)

var Reporter = goerrorcollector.GoErrorCollector{}

func Report(component string, severity int, message string, stacktrace string) {
	Reporter.Report(goerrorcollector.ErrorMessage{
		Component:  component,
		Message:    message,
		Stacktrace: stacktrace,
		Time:       time.Now().UnixMilli(),
		Severity:   severity,
	})
}

func Start(url string, exchange string, topic string) {
	Reporter.Start(goerrorcollector.WithReceiverRabbitMQ(goerrorcollector.RecieverMqOption{
		RabbitMQUrl:  url,
		ExchangeName: exchange,
		Topic:        topic,
	}))
}
