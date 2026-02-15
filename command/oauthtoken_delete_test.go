package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

type mockOAuthTokenDeleter struct {
	lastID     string
	deleteFunc func(ctx context.Context, oAuthTokenID string) error
}

func (m *mockOAuthTokenDeleter) Delete(ctx context.Context, oAuthTokenID string) error {
	m.lastID = oAuthTokenID
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, oAuthTokenID)
	}
	return nil
}

func TestOAuthTokenDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthTokenDeleteCommand{Meta: newTestMeta(ui)}

	code := cmd.Run(nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}
}

func TestOAuthTokenDeleteSuccessWithForce(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOAuthTokenDeleter{}
	cmd := &OAuthTokenDeleteCommand{
		Meta:          testMeta(t, ui),
		oauthTokenSvc: svc,
	}

	code := cmd.Run([]string{"-id=ot-123", "-force"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; err=%q", code, ui.ErrorWriter.String())
	}
	if svc.lastID != "ot-123" {
		t.Fatalf("expected delete call with ot-123, got %q", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success output, got %q", ui.OutputWriter.String())
	}
}

func TestOAuthTokenDeleteDeleteError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOAuthTokenDeleter{
		deleteFunc: func(ctx context.Context, oAuthTokenID string) error {
			return errors.New("boom")
		},
	}
	cmd := &OAuthTokenDeleteCommand{
		Meta:          testMeta(t, ui),
		oauthTokenSvc: svc,
	}

	code := cmd.Run([]string{"-id=ot-123", "-y"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}
