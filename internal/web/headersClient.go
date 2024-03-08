package web

import (
	"context"
	"io"
)

type ClientWithHeaders struct {
	client  *Client
	headers map[string]string
}

func (c *Client) WithHeaders(headers map[string]string) *ClientWithHeaders {
	return &ClientWithHeaders{
		client:  c,
		headers: headers,
	}
}

func (c *ClientWithHeaders) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	return c.client.DownloadWithHeaders(ctx, path, c.headers)
}
