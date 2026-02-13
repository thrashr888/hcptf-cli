package command

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func TestConfigVersionUploadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresPath(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=cv-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-path") {
		t.Fatalf("expected path error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-path=/tmp"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresEmptyPath(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=cv-123", "-path="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-path") {
		t.Fatalf("expected path error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresPathExists(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=cv-123", "-path=/nonexistent/path/that/does/not/exist"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "does not exist") {
		t.Fatalf("expected path existence error, got %q", out)
	}
}

func TestConfigVersionUploadSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	temp := t.TempDir()
	path := temp + "/config.txt"
	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}

	reader := &mockConfigVersionReadService{
		response: &tfe.ConfigurationVersion{
			ID:        "cv-1",
			UploadURL: "https://upload.example.com",
		},
	}
	uploader := &mockConfigVersionUploader{}

	cmd := &ConfigVersionUploadCommand{
		Meta:         newTestMeta(ui),
		configVerSvc: reader,
		uploadSvc:    uploader,
	}

	code := cmd.Run([]string{"-id=cv-1", "-path=" + path})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if reader.lastID != "cv-1" {
		t.Fatalf("expected read called for cv-1, got %q", reader.lastID)
	}
	if uploader.lastURL != "https://upload.example.com" || uploader.lastPath != path {
		t.Fatalf("expected upload called with url=%q path=%q, got url=%q path=%q",
			"https://upload.example.com", path, uploader.lastURL, uploader.lastPath)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Successfully uploaded configuration to version: cv-1") {
		t.Fatalf("expected success output, got %q", ui.OutputWriter.String())
	}
}

func TestConfigVersionUploadReadError(t *testing.T) {
	ui := cli.NewMockUi()
	temp := t.TempDir()
	path := temp + "/config.txt"
	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}

	reader := &mockConfigVersionReadService{
		err: errors.New("read failed"),
	}
	uploader := &mockConfigVersionUploader{}

	cmd := &ConfigVersionUploadCommand{
		Meta:         newTestMeta(ui),
		configVerSvc: reader,
		uploadSvc:    uploader,
	}

	code := cmd.Run([]string{"-id=cv-1", "-path=" + path})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error reading configuration version: read failed") {
		t.Fatalf("expected read error output, got %q", ui.ErrorWriter.String())
	}
	if uploader.lastURL != "" {
		t.Fatal("upload should not be called when read fails")
	}
}

func TestConfigVersionUploadMissingUploadURL(t *testing.T) {
	ui := cli.NewMockUi()
	temp := t.TempDir()
	path := temp + "/config.txt"
	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}

	reader := &mockConfigVersionReadService{
		response: &tfe.ConfigurationVersion{
			ID: "cv-1",
		},
	}
	uploader := &mockConfigVersionUploader{}

	cmd := &ConfigVersionUploadCommand{
		Meta:         newTestMeta(ui),
		configVerSvc: reader,
		uploadSvc:    uploader,
	}

	code := cmd.Run([]string{"-id=cv-1", "-path=" + path})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "does not have an upload URL") {
		t.Fatalf("expected missing upload URL output, got %q", ui.ErrorWriter.String())
	}
	if uploader.lastURL != "" {
		t.Fatal("upload should not be called without upload URL")
	}
}

func TestConfigVersionUploadUploadError(t *testing.T) {
	ui := cli.NewMockUi()
	temp := t.TempDir()
	path := temp + "/config.txt"
	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}

	reader := &mockConfigVersionReadService{
		response: &tfe.ConfigurationVersion{
			ID:        "cv-1",
			UploadURL: "https://upload.example.com",
		},
	}
	uploader := &mockConfigVersionUploader{
		err: errors.New("upload failed"),
	}

	cmd := &ConfigVersionUploadCommand{
		Meta:         newTestMeta(ui),
		configVerSvc: reader,
		uploadSvc:    uploader,
	}

	code := cmd.Run([]string{"-id=cv-1", "-path=" + path})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error uploading configuration: upload failed") {
		t.Fatalf("expected upload error output, got %q", ui.ErrorWriter.String())
	}
}

type mockConfigVersionUploader struct {
	lastURL  string
	lastPath string
	err      error
}

func (m *mockConfigVersionUploader) Upload(_ context.Context, url, path string) error {
	m.lastURL = url
	m.lastPath = path
	return m.err
}
