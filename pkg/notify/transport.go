package notify

import (
	"context"
	"fmt"
	c "github.com/ncostamagna/go_http_client/client"
	"net/http"
	"net/url"
	"time"
)

type (
	DataResponse struct {
		Success string   `json:"success"`
		Error   string `json: "error"`
	}

	Transport interface {
		Push(ctx context.Context, title, message, urlNotify string) error
	}

	clientHTTP struct {
		client c.Transport
		token string
	}
)

func NewHttpClient(baseURL, token string) Transport {
	header := http.Header{}

	return &clientHTTP{
		client: c.New(header, baseURL, 5000*time.Millisecond, true),
		token: token,
	}
}

func (c *clientHTTP) Push(ctx context.Context, title, message, urlNotify string) error {

	dataResponse := DataResponse{}

	u := url.URL{}
	u.Path = "/api/message"

	if c.token == "" || title == "" || message == "" {
		return fmt.Errorf("token, title and message are required")
	}

	query := u.Query()
	query.Set("k", c.token)
	query.Set("t", title)
	query.Set("c", message)

	if urlNotify != "" {
		query.Set("u", urlNotify)
	}
	u.RawQuery = query.Encode()

	reps := c.client.Get(u.String())

	if reps.Err != nil {
		return reps.Err
	}

	if err := reps.FillUp(&dataResponse); err != nil {
		return fmt.Errorf("%s", reps)
	}

	if reps.StatusCode > 299 {
		return fmt.Errorf("%s", dataResponse.Error)
	}

	if dataResponse.Success != "1" {
		return fmt.Errorf("%s", dataResponse.Error)
	}

	return nil
}