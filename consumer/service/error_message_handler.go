package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/latifrons/goerrorcollector"
	"github.com/rabbitmq/amqp091-go"
	"net/url"
	"time"
)
import "github.com/jellydator/ttlcache/v3"

type ErrorMessageHandlerImpl struct {
	DiscordWebhookUrl    string
	ProxyUrl             string
	duplicateErrorCache  *ttlcache.Cache[string, int64]
	discordRequestSender *DiscordRequestSender
}

func (e *ErrorMessageHandlerImpl) InitDefault() (err error) {
	e.duplicateErrorCache = ttlcache.New[string, int64](
		ttlcache.WithTTL[string, int64](time.Minute),
		//ttlcache.WithDisableTouchOnHit[string, int64](true)
	)
	go e.duplicateErrorCache.Start()

	var p *url.URL
	if e.ProxyUrl != "" {
		p, err = url.Parse(e.ProxyUrl)
		if err != nil {
			return
		}
	}
	e.discordRequestSender = NewDiscordRequestSender(p)
	return
}

func (e *ErrorMessageHandlerImpl) HandleErrorMessage(delivery amqp091.Delivery) (err error) {
	var errorMessage goerrorcollector.ErrorMessage
	err = json.Unmarshal(delivery.Body, &errorMessage)
	if err != nil {
		return
	}

	// send to discord
	err = e.SendDiscordMessage(errorMessage)
	if err != nil {
		return
	}
	return
}

func (e *ErrorMessageHandlerImpl) SendDiscordMessage(message goerrorcollector.ErrorMessage) (err error) {
	author := fmt.Sprintf("[%d] %s", message.Severity, message.Component)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	err = e.discordRequestSender.SendMessage(ctx, e.DiscordWebhookUrl, Message{
		Username: &author,
		Content:  &message.Message,
	})
	if err != nil {
		return err
	}
	return
}
