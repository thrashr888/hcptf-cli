package command

import (
	"testing"
	"time"
)

func TestNewHTTPClientTimeout(t *testing.T) {
	client := newHTTPClient()
	if client.Timeout != 30*time.Second {
		t.Fatalf("expected timeout to be 30s, got %s", client.Timeout)
	}
}
