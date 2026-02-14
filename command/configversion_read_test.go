package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newConfigVersionReadCommand(ui cli.Ui, svc configVersionReader) *ConfigVersionReadCommand {
	return &ConfigVersionReadCommand{
		Meta:         newTestMeta(ui),
		configVerSvc: svc,
	}
}

func TestConfigVersionReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newConfigVersionReadCommand(ui, &mockConfigVersionReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestConfigVersionReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockConfigVersionReadService{err: errors.New("boom")}
	cmd := newConfigVersionReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=cv-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "cv-1" {
		t.Fatalf("expected config version id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestConfigVersionReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockConfigVersionReadService{response: &tfe.ConfigurationVersion{
		ID:          "cv-1",
		Status:      tfe.ConfigurationUploaded,
		Source:      tfe.ConfigurationSourceAPI,
		Speculative: true,
		Error:       "",
	}}
	cmd := newConfigVersionReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=cv-1", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "cv-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestConfigVersionReadRunID(t *testing.T) {
	ui := cli.NewMockUi()
	configSvc := &mockConfigVersionReadService{response: &tfe.ConfigurationVersion{
		ID:     "cv-123",
		Status: tfe.ConfigurationUploaded,
	}}
	runSvc := &mockRunReadService{response: &tfe.Run{
		ID: "run-123",
		ConfigurationVersion: &tfe.ConfigurationVersion{
			ID: "cv-123",
		},
	}}
	cmd := &ConfigVersionReadCommand{
		Meta:         newTestMeta(ui),
		configVerSvc: configSvc,
		runSvc:       runSvc,
	}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-run-id=run-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "cv-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
	if runSvc.lastRun != "run-123" {
		t.Fatalf("expected run id recorded, got: %s", runSvc.lastRun)
	}
	if configSvc.lastID != "cv-123" {
		t.Fatalf("expected config version read id recorded, got: %s", configSvc.lastID)
	}
}

func TestConfigVersionReadRunReadFailure(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{err: errors.New("run failed")}
	cmd := &ConfigVersionReadCommand{
		Meta:         newTestMeta(ui),
		runSvc:       runSvc,
		configVerSvc: &mockConfigVersionReadService{},
	}

	if code := cmd.Run([]string{"-run-id=run-999"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if runSvc.lastRun != "run-999" {
		t.Fatalf("expected run id recorded, got: %s", runSvc.lastRun)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error reading run: run failed") {
		t.Fatalf("expected run read error output, got: %s", ui.ErrorWriter.String())
	}
}

func TestConfigVersionReadRunWithNoConfigVersion(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{response: &tfe.Run{
		ID: "run-123",
	}}
	cmd := &ConfigVersionReadCommand{
		Meta:         newTestMeta(ui),
		runSvc:       runSvc,
		configVerSvc: &mockConfigVersionReadService{},
	}

	if code := cmd.Run([]string{"-id=run-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error: run has no configuration version") {
		t.Fatalf("expected no configuration version error, got: %s", ui.ErrorWriter.String())
	}
}
