package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/inhandnet/inconnect-cli/internal/debug"
)

type HTTPError struct {
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, string(e.Body))
}

type APIClient struct {
	client  *resty.Client
	baseURL string
	verbose int
}

func NewAPIClient(baseURL string, transport http.RoundTripper, verbose int) *APIClient {
	c := resty.New()
	c.SetTransport(transport)
	c.SetBaseURL(baseURL)
	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")
	return &APIClient{client: c, baseURL: baseURL, verbose: verbose}
}

func (c *APIClient) BaseURL() string { return c.baseURL }

func (c *APIClient) Get(path string, query url.Values) ([]byte, error) {
	if c.verbose > 0 {
		if query == nil {
			query = url.Values{}
		}
		if query.Get("verbose") == "" {
			query.Set("verbose", strconv.Itoa(c.verbose))
		}
	}
	cleanEmpty(query)
	req := c.client.R().SetQueryParamsFromValues(query)
	return c.execute(req, resty.MethodGet, path)
}

func (c *APIClient) Post(path string, body interface{}) ([]byte, error) {
	req := c.client.R().SetBody(body)
	return c.execute(req, resty.MethodPost, path)
}

func (c *APIClient) Put(path string, body interface{}) ([]byte, error) {
	req := c.client.R().SetBody(body)
	return c.execute(req, resty.MethodPut, path)
}

func (c *APIClient) Delete(path string) ([]byte, error) {
	req := c.client.R()
	return c.execute(req, resty.MethodDelete, path)
}

func (c *APIClient) Do(method, path string, query url.Values, body interface{}) ([]byte, error) {
	cleanEmpty(query)
	req := c.client.R().SetQueryParamsFromValues(query)
	if body != nil {
		req.SetBody(body)
	}
	return c.execute(req, method, path)
}

func (c *APIClient) execute(req *resty.Request, method, path string) ([]byte, error) {
	resp, err := req.Execute(method, path)
	if err != nil {
		return nil, err
	}

	body := resp.Body()
	debug.Log("response %d: %s", resp.StatusCode(), string(body))

	if resp.StatusCode() >= 400 {
		return body, &HTTPError{StatusCode: resp.StatusCode(), Body: body}
	}
	if err := BodyError(body); err != nil {
		return body, err
	}
	return body, nil
}

func cleanEmpty(q url.Values) {
	for k, v := range q {
		if len(v) == 1 && v[0] == "" {
			delete(q, k)
		}
	}
}
