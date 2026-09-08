package remote

import (
	"bytes"
	"context"
	"net/http"
)

type HTTPClient struct {
	client  *http.Client
	headers map[string]string
}

func NewHTTPClient(
	client *http.Client,
	headers map[string]string,
) *HTTPClient {
	if client == nil {
		client = http.DefaultClient
	}

	return &HTTPClient{
		client:  client,
		headers: headers,
	}
}

func (c *HTTPClient) Do(
	ctx context.Context,
	method string,
	url string,
	body []byte,
	headers map[string]string,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	// Apply default headers.
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Apply request-specific headers.
	// These override default headers if the same key exists.
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.client.Do(req)
}
