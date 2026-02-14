package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newAccountShowCommand(ui cli.Ui, svc accountReader) *AccountShowCommand {
	return &AccountShowCommand{
		Meta:       newTestMeta(ui),
		accountSvc: svc,
	}
}

func TestAccountShowHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAccountReadService{err: errors.New("unauthorized")}
	cmd := newAccountShowCommand(ui, svc)

	code := cmd.Run(nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "unauthorized") {
		t.Fatalf("expected error output")
	}
}

func TestAccountShowSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAccountReadService{response: &tfe.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Username: "testuser",
	}}
	cmd := newAccountShowCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run(nil)
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	_ = output // output goes to stdout via formatter
}

func TestAccountShowWithTwoFactor(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAccountReadService{response: &tfe.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Username: "testuser",
		TwoFactor: &tfe.TwoFactor{
			Enabled:  true,
			Verified: true,
		},
	}}
	cmd := newAccountShowCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run(nil)
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	_ = output
}

func TestAccountShowHelp(t *testing.T) {
	cmd := &AccountShowCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "account show") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestAccountShowSynopsis(t *testing.T) {
	cmd := &AccountShowCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
