package client

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

type ClientIface interface {
	GetDetail(ctx context.Context, song string, group string) (models.Song, error)
}

type Client struct {
	client   *http.Client
	basePath url.URL
}

func New(host string, port string) *Client {
	return &Client{
		client: &http.Client{},
		basePath: url.URL{
			Scheme: "http",
			Host:   host + ":" + port,
			Path:   "/info",
		},
	}
}

func (c *Client) GetDetail(ctx context.Context, song string, group string) (models.Song, error) {
	logger.LogUse(ctx).Debug("Client.GetDetail", "input", slog.String("song", song), slog.String("group", group))

	p := url.Values{}
	p.Add("song", song)
	p.Add("group", group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.basePath.String(), nil)
	if err != nil {
		return models.Song{}, err
	}

	req.URL.RawQuery = p.Encode()
	logger.LogUse(ctx).Info("Do Request", slog.String("url", req.URL.String()))

	resp, err := c.client.Do(req)
	if err != nil {
		return models.Song{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Song{}, err
	}

	var res models.Song

	if err := json.Unmarshal(body, &res); err != nil {
		return models.Song{}, err
	}

	resp.Body.Close()

	logger.LogUse(ctx).Debug("Result", slog.Any("song", res.AsLogValue()))

	return res, nil
}
