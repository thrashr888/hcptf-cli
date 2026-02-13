package command

import (
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newUserTokenListCommand(ui cli.Ui, userSvc userReader, tokenSvc userTokenLister) *UserTokenListCommand {
	return &UserTokenListCommand{
		Meta:         newTestMeta(ui),
		userSvc:      userSvc,
		userTokenSvc: tokenSvc,
	}
}

func TestUserTokenListHandlesUserError(t *testing.T) {
	ui := cli.NewMockUi()
	userSvc := &mockUserReadService{err: errors.New("unauthorized")}
	tokenSvc := &mockUserTokenListService{}
	cmd := newUserTokenListCommand(ui, userSvc, tokenSvc)

	code := cmd.Run(nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "unauthorized") {
		t.Fatalf("expected error output")
	}
}

func TestUserTokenListHandlesTokenError(t *testing.T) {
	ui := cli.NewMockUi()
	userSvc := &mockUserReadService{response: &tfe.User{ID: "user-1"}}
	tokenSvc := &mockUserTokenListService{err: errors.New("forbidden")}
	cmd := newUserTokenListCommand(ui, userSvc, tokenSvc)

	code := cmd.Run(nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if tokenSvc.lastID != "user-1" {
		t.Fatalf("expected user ID user-1, got %s", tokenSvc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "forbidden") {
		t.Fatalf("expected error output")
	}
}

func TestUserTokenListEmpty(t *testing.T) {
	ui := cli.NewMockUi()
	userSvc := &mockUserReadService{response: &tfe.User{ID: "user-1"}}
	tokenSvc := &mockUserTokenListService{response: &tfe.UserTokenList{Items: []*tfe.UserToken{}}}
	cmd := newUserTokenListCommand(ui, userSvc, tokenSvc)

	code := cmd.Run(nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "No user tokens") {
		t.Fatalf("expected no tokens message")
	}
}

func TestUserTokenListSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	userSvc := &mockUserReadService{response: &tfe.User{ID: "user-1"}}
	tokenSvc := &mockUserTokenListService{response: &tfe.UserTokenList{
		Items: []*tfe.UserToken{
			{
				ID:          "at-1",
				Description: "CI token",
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}}
	cmd := newUserTokenListCommand(ui, userSvc, tokenSvc)

	_, code := captureStdout(t, func() int {
		return cmd.Run(nil)
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
}

func TestUserTokenListHelp(t *testing.T) {
	cmd := &UserTokenListCommand{}
	if !strings.Contains(cmd.Help(), "usertoken list") {
		t.Fatal("expected help text")
	}
}

func TestUserTokenListSynopsis(t *testing.T) {
	cmd := &UserTokenListCommand{}
	if cmd.Synopsis() == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
