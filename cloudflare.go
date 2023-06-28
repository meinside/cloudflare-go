package cfgo

import (
	"net"
	"net/http"
	"time"
)

const (
	timeoutSeconds = 10
)

// CloudflareClient struct
type CloudflareClient struct {
	Email  string
	APIKey string

	httpClient *http.Client

	Verbose bool
}

// NewCloudflareClient returns a new cloudflare API client with given api key.
func NewCloudflareClient(email, apiKey string) *CloudflareClient {
	return &CloudflareClient{
		Email:  email,
		APIKey: apiKey,

		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   timeoutSeconds * time.Second,
					KeepAlive: timeoutSeconds * time.Second,
				}).DialContext,
				IdleConnTimeout:       timeoutSeconds * time.Second,
				TLSHandshakeTimeout:   timeoutSeconds * time.Second,
				ResponseHeaderTimeout: timeoutSeconds * time.Second,
				ExpectContinueTimeout: timeoutSeconds * time.Second,
			},
		},
	}
}
