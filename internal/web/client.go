package web

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/shushard/ContestCourier/internal/web/config"
	"github.com/shushard/ContestCourier/pkg/utils"
)

type Client struct {
	config     *config.Config
	httpClient *http.Client
	logger     *slog.Logger
}

func New(conf *config.Config, logger *slog.Logger) *Client {
	return &Client{
		config: conf,
		httpClient: &http.Client{
			Timeout: conf.Timeout,
		},
		logger: logger,
	}
}

func (c *Client) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	return c.DownloadWithHeaders(ctx, path, nil)
}

func (c *Client) DownloadWithHeaders(
	ctx context.Context,
	path string,
	headers map[string]string,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request with path=%s: %w", path, err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	var data io.ReadCloser
	err = utils.WithRetries(ctx, uint(c.config.Retries), c.config.RetryDelay, func() (err error) {
		var resp *http.Response
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("can't do web request: %w", err)
		}

		setStatusCode(path, resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			var body []byte
			body, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			return newHTTPError(resp.StatusCode, body, err)
		}

		data = resp.Body
		return nil
	}, c.logger)
	if err != nil {
		return nil, fmt.Errorf("web gave up downloading: %w", err)
	}

	return data, nil
}
