package command

import (
	"io"
	"os"
	"testing"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/mitchellh/cli"
)

func newTestMeta(ui cli.Ui) Meta {
	return Meta{Ui: ui, client: &client.Client{}}
}

func captureStdout(t *testing.T, fn func() int) (string, int) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()
	code := fn()
	w.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	return string(data), code
}
