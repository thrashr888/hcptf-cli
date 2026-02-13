package command

import (
	"net/http"
	"time"
)

const defaultHTTPClientTimeout = 30 * time.Second

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultHTTPClientTimeout,
	}
}
