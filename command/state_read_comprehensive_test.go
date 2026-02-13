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

func newStateReadCommand(ui cli.Ui, svc stateVersionReader) *StateReadCommand {
	return &StateReadCommand{
		Meta:     newTestMeta(ui),
		stateSvc: svc,
	}
}

func TestStateReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newStateReadCommand(ui, &mockStateVersionReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestStateReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockStateVersionReadService{err: errors.New("boom")}
	cmd := newStateReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=sv-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "sv-123" {
		t.Fatalf("unexpected state version id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestStateReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockStateVersionReadService{
		response: &tfe.StateVersion{
			ID:                "sv-123",
			Serial:            5,
			CreatedAt:         time.Unix(0, 0),
			DownloadURL:       "https://example.com/state",
			Resources:         []*tfe.StateVersionResources{{Name: "resource1"}, {Name: "resource2"}},
			ResourcesProcessed: true,
			StateVersion:      4,
			TerraformVersion:  "1.5.0",
		},
	}
	cmd := newStateReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=sv-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "sv-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
