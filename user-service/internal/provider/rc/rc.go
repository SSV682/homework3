package rc

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

const defaultRequestTimeout = 5 * time.Second

type httpClient struct {
	client   *fasthttp.Client
	endpoint string
	baseURL  string
}

func (c *httpClient) RequestURL(userID string) string {
	return c.baseURL + c.endpoint + userID
}

func NewClientProvider(baseURL, endpoint string) *httpClient {
	return &httpClient{
		baseURL:  baseURL,
		endpoint: endpoint,
		client:   &fasthttp.Client{},
	}
}

func (c *httpClient) CreateAccount(userID string) error {
	req := fasthttp.AcquireRequest()

	log.Infof("url: %s", c.RequestURL(userID))
	req.SetRequestURI(c.RequestURL(userID))
	req.Header.SetMethod(fasthttp.MethodPost)
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.client.DoTimeout(req, resp, defaultRequestTimeout); err != nil {
		return fmt.Errorf("http request to create account: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("bad http response status %d", resp.StatusCode())
	}

	return nil
}
