package infra

import (
	"fmt"
	"net/url"

	"github.com/minio/madmin-go/v4"
)

type MinIOConfig struct {
	URL      string
	RootUser string
	Password string
}

func NewMinIOClient(cfg MinIOConfig) (*madmin.AdminClient, error) {
	endpoint, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid MinIO URL: %w", err)
	}

	useSSL := endpoint.Scheme == "https"
	host := endpoint.Host

	client, err := madmin.New(host, cfg.RootUser, cfg.Password, useSSL)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO admin client: %w", err)
	}

	return client, nil
}
