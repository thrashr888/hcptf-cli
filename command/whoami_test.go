package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newWhoAmICommand(ui cli.Ui, svc accountReader) *WhoAmICommand {
	return &WhoAmICommand{
		Meta:       newTestMeta(ui),
		accountSvc: svc,
	}
}

func TestWhoAmICommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAccountReadService{err: errors.New("unauthorized")}
	cmd := newWhoAmICommand(ui, svc)

	code := cmd.Run(nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "unauthorized") {
		t.Fatalf("expected error output")
	}
}

func TestWhoAmICommandSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAccountReadService{response: &tfe.User{
		ID:               "user-1",
		Email:            "test@example.com",
		Username:         "testuser",
		IsServiceAccount: false,
	}}
	cmd := newWhoAmICommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run(nil)
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(output, "testuser") || !strings.Contains(output, "test@example.com") {
		t.Fatalf("expected output to include user identity, got: %s", output)
	}
}

func TestWhoAmICommandJSONOutput(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAccountReadService{response: &tfe.User{
		ID:               "user-1",
		Email:            "test@example.com",
		Username:         "testuser",
		IsServiceAccount: true,
	}}
	cmd := newWhoAmICommand(ui, svc)
	cmd.format = "json"

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if got, ok := payload["ID"]; !ok || got != "user-1" {
		t.Fatalf("expected ID=user-1, got: %v", payload["ID"])
	}
	if got, ok := payload["IsServiceAccount"]; !ok || got != true {
		t.Fatalf("expected IsServiceAccount=true, got: %v", payload["IsServiceAccount"])
	}
}

func TestWhoAmICommandHelp(t *testing.T) {
	cmd := &WhoAmICommand{}
	if !strings.Contains(cmd.Help(), "hcptf whoami") {
		t.Fatalf("expected help text, got: %s", cmd.Help())
	}
}

func TestWhoAmICommandSynopsis(t *testing.T) {
	cmd := &WhoAmICommand{}
	if cmd.Synopsis() == "" {
		t.Fatal("expected non-empty synopsis")
	}
}

func TestWhoAmICommandFactory(t *testing.T) {
	meta := newTestMeta(cli.NewMockUi())
	commands := Commands(&meta)

	factory, ok := commands["whoami"]
	if !ok {
		t.Fatal("expected whoami command in command map")
	}
	cmd, err := factory()
	if err != nil {
		t.Fatalf("whoami factory returned error: %v", err)
	}
	if _, ok := cmd.(*WhoAmICommand); !ok {
		t.Fatalf("expected *WhoAmICommand, got %T", cmd)
	}
}
