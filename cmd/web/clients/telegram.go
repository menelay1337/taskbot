package telegram

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"encoding/json"
	
	"taskbot/internal/errors"
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

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return errors.WrapIfErr("can't send message:", err)
	}

	return nil
}



func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse 

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}


}

const (
	getUpdatesMethod = "getUpdates",
	sendMessageMethod = "sendMessage",
)

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = errors.WrapIfErr("can't do request", err) }()

	u := url.URL {
		Scheme: "https",
		Host:	c.host,
		Path:	path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet,  u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

