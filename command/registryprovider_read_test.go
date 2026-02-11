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

func newRegistryProviderReadCommand(ui cli.Ui, svc registryProviderReader) *RegistryProviderReadCommand {
	return &RegistryProviderReadCommand{
		Meta:                newTestMeta(ui),
		registryProviderSvc: svc,
	}
}

func TestRegistryProviderReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryProviderReadCommand(ui, &mockRegistryProviderReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 name")
	}
}

func TestRegistryProviderReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderReadService{err: errors.New("boom")}
	cmd := newRegistryProviderReadCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=aws"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID.Name != "aws" {
		t.Fatalf("expected provider name recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRegistryProviderReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderReadService{response: &tfe.RegistryProvider{ID: "prov-1", Name: "aws", Namespace: "hashicorp", RegistryName: tfe.PrivateRegistry, CreatedAt: time.Now().Format(time.RFC3339)}}
	cmd := newRegistryProviderReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=aws", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["Name"] != "aws" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
