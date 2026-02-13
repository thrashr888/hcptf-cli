package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newCommentListCommand(ui cli.Ui, svc commentLister) *CommentListCommand {
	return &CommentListCommand{
		Meta:       newTestMeta(ui),
		commentSvc: svc,
	}
}

func TestCommentListRequiresRunID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newCommentListCommand(ui, &mockCommentListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-run-id") {
		t.Fatalf("expected run-id error")
	}
}

func TestCommentListHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockCommentListService{err: errors.New("forbidden")}
	cmd := newCommentListCommand(ui, svc)

	code := cmd.Run([]string{"-run-id=run-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastID != "run-123" {
		t.Fatalf("expected run ID run-123, got %s", svc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "forbidden") {
		t.Fatalf("expected error output")
	}
}

func TestCommentListEmpty(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockCommentListService{response: &tfe.CommentList{Items: []*tfe.Comment{}}}
	cmd := newCommentListCommand(ui, svc)

	code := cmd.Run([]string{"-run-id=run-123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "No comments") {
		t.Fatalf("expected no comments message")
	}
}

func TestCommentListSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockCommentListService{response: &tfe.CommentList{
		Items: []*tfe.Comment{
			{ID: "comment-1", Body: "LGTM"},
			{ID: "comment-2", Body: "Please fix the tests"},
		},
	}}
	cmd := newCommentListCommand(ui, svc)

	_, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-run-id=run-123"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
}

func TestCommentListHelp(t *testing.T) {
	cmd := &CommentListCommand{}
	if !strings.Contains(cmd.Help(), "comment list") {
		t.Fatal("expected help text")
	}
}

func TestCommentListSynopsis(t *testing.T) {
	cmd := &CommentListCommand{}
	if cmd.Synopsis() == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
