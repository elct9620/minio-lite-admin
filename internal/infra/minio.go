package infra

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/madmin-go/v4"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	client, err := madmin.NewWithOptions(host, &madmin.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.Password, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO admin client: %w", err)
	}

	return client, nil
}

// NewMinIOClientWithTimeout creates a MinIO admin client with custom timeout for testing
func NewMinIOClientWithTimeout(cfg MinIOConfig, timeout time.Duration) (*madmin.AdminClient, error) {
	endpoint, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid MinIO URL: %w", err)
	}

	useSSL := endpoint.Scheme == "https"
	host := endpoint.Host

	// Create custom HTTP transport with short timeout
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout:   timeout,
		ResponseHeaderTimeout: timeout,
	}

	// Create client with custom options
	opts := &madmin.Options{
		Creds:     credentials.NewStaticV4(cfg.RootUser, cfg.Password, ""),
		Secure:    useSSL,
		Transport: transport,
	}

	client, err := madmin.NewWithOptions(host, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO admin client with timeout: %w", err)
	}

	return client, nil
}
