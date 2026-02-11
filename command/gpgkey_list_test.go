package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newGPGKeyListCommand(ui cli.Ui, svc gpgKeyLister) *GPGKeyListCommand {
	return &GPGKeyListCommand{
		Meta:      newTestMeta(ui),
		gpgKeySvc: svc,
	}
}

func TestGPGKeyListRequiresNamespace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyListCommand(ui, &mockGPGKeyListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-namespace") {
		t.Fatalf("expected namespace error")
	}
}

func TestGPGKeyListHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyListService{err: errors.New("boom")}
	cmd := newGPGKeyListCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOptions.Namespaces[0] != "org" {
		t.Fatalf("expected namespace in options")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestGPGKeyListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyListService{response: &tfe.GPGKeyList{Items: []*tfe.GPGKey{{ID: "key-1", KeyID: "abc", Namespace: "org"}}}}
	cmd := newGPGKeyListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-namespace=org", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if rows[0]["Key ID"] != "abc" {
		t.Fatalf("unexpected row: %#v", rows)
	}
}
