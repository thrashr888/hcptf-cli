package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newGPGKeyReadCommand(ui cli.Ui, svc gpgKeyReader) *GPGKeyReadCommand {
	return &GPGKeyReadCommand{
		Meta:      newTestMeta(ui),
		gpgKeySvc: svc,
	}
}

func TestGPGKeyReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyReadCommand(ui, &mockGPGKeyReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 namespace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-namespace") {
		t.Fatalf("expected namespace error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-namespace=org"}); code != 1 {
		t.Fatalf("expected exit 1 key")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-key-id") {
		t.Fatalf("expected key error")
	}
}

func TestGPGKeyReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyReadService{err: errors.New("boom")}
	cmd := newGPGKeyReadCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org", "-key-id=abc"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID.Namespace != "org" || svc.lastID.KeyID != "abc" {
		t.Fatalf("unexpected key id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestGPGKeyReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyReadService{response: &tfe.GPGKey{ID: "key-1", KeyID: "abc", Namespace: "org", CreatedAt: time.Unix(0, 0)}}
	cmd := newGPGKeyReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-namespace=org", "-key-id=abc", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "key-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
