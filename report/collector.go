package report

import (
	goerrorcollector "github.com/latifrons/go-error-collector"
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
func ReportError(component string, severity int, message string, err error) {
	if err != nil {
		Reporter.Report(goerrorcollector.ErrorMessage{
			Component:  component,
			Message:    message + " " + err.Error(),
			Stacktrace: "",
			Time:       time.Now().UnixMilli(),
			Severity:   severity,
		})
	} else {
		Reporter.Report(goerrorcollector.ErrorMessage{
			Component:  component,
			Message:    message,
			Stacktrace: "",
			Time:       time.Now().UnixMilli(),
			Severity:   severity,
		})
	}
}

func Start(url string, exchange string, topic string) {
	Reporter.Start(goerrorcollector.WithReceiverRabbitMQ(goerrorcollector.ReceiverMqOption{
		RabbitMQUrl:  url,
		ExchangeName: exchange,
		Topic:        topic,
	}))
}
