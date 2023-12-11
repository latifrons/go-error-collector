package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Message struct {
	Username        *string          `json:"username,omitempty"`
	AvatarUrl       *string          `json:"avatar_url,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Embeds          *[]Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
}

type Embed struct {
	Title       *string    `json:"title,omitempty"`
	Url         *string    `json:"url,omitempty"`
	Description *string    `json:"description,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      *[]Field   `json:"fields,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
}

type Author struct {
	Name    *string `json:"name,omitempty"`
	Url     *string `json:"url,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   *string `json:"name,omitempty"`
	Value  *string `json:"value,omitempty"`
	Inline *bool   `json:"inline,omitempty"`
}

type Thumbnail struct {
	Url *string `json:"url,omitempty"`
}

type Image struct {
	Url *string `json:"url,omitempty"`
}

type Footer struct {
	Text    *string `json:"text,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}

type AllowedMentions struct {
	Parse *[]string `json:"parse,omitempty"`
	Users *[]string `json:"users,omitempty"`
	Roles *[]string `json:"roles,omitempty"`
}

type DiscordRequestSender struct {
	ProxyUrl   *url.URL
	httpClient http.Client
}

func NewDiscordRequestSender(proxyUrl *url.URL) *DiscordRequestSender {
	return &DiscordRequestSender{
		ProxyUrl: proxyUrl,
		httpClient: http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		},
	}
}

func (s *DiscordRequestSender) SendMessage(ctx context.Context, requestUrl string, message Message) (err error) {
	payload := new(bytes.Buffer)

	err = json.NewEncoder(payload).Encode(message)
	if err != nil {
		return err
	}

	// Build the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}
