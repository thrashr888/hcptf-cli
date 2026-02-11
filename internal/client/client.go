package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/config"
)

// Client wraps the go-tfe client with CLI-specific functionality
type Client struct {
	*tfe.Client
	config  *config.Config
	address string
}

// New creates a new API client
func New(cfg *config.Config) (*Client, error) {
	address := config.GetAddress()

	// Parse the address to get hostname for token lookup
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("invalid API address: %w", err)
	}

	token := cfg.GetToken(u.Hostname())
	if token == "" {
		return nil, fmt.Errorf("no authentication token found. Set HCPTF_TOKEN environment variable or configure credentials in ~/.hcptfrc")
	}

	// Create TFE client configuration
	tfeConfig := &tfe.Config{
		Address:    address,
		Token:      token,
		HTTPClient: http.DefaultClient,
	}

	// Create the TFE client
	tfeClient, err := tfe.NewClient(tfeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &Client{
		Client:  tfeClient,
		config:  cfg,
		address: address,
	}, nil
}

// GetAddress returns the API address being used
func (c *Client) GetAddress() string {
	return c.address
}

// Context returns a background context
// Commands can override this if they need cancellation or timeouts
func (c *Client) Context() context.Context {
	return context.Background()
}
