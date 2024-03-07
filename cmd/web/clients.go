package telegram

import (
	"net/http"
)

type Client struct {
	host	 string // 
	basePath string // tg-bot.com/bot<token>
	client http.Client
}

func New(host string, token string) {
	return Client {
		host:		host,
		basePath:   newBasePath(token),
		client:		http.Client{},
	}

}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates() {
	 	
}

func (c *Client) SendMessage() {

}
