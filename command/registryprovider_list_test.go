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

func newRegistryProviderListCommand(ui cli.Ui, svc registryProviderLister) *RegistryProviderListCommand {
	return &RegistryProviderListCommand{
		Meta:                newTestMeta(ui),
		registryProviderSvc: svc,
	}
}

func TestRegistryProviderListRequiresOrg(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryProviderListCommand(ui, &mockRegistryProviderListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}
}

func TestRegistryProviderListHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderListService{err: errors.New("boom")}
	cmd := newRegistryProviderListCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected org recorded")
	}
	if svc.lastOpts == nil || svc.lastOpts.ListOptions.PageSize != 100 {
		t.Fatalf("expected list options set")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error message")
	}
}

func TestRegistryProviderListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderListService{response: &tfe.RegistryProviderList{Items: []*tfe.RegistryProvider{{
		ID:           "prov-1",
		Name:         "aws",
		Namespace:    "hashicorp",
		RegistryName: tfe.PrivateRegistry,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}}}}
	cmd := newRegistryProviderListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if rows[0]["Name"] != "aws" {
		t.Fatalf("unexpected row data: %#v", rows)
	}
}
