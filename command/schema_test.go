package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestParseFlagSchemaFromHelp(t *testing.T) {
	help := `
Usage: hcptf workspace create [options]

Options:

  -organization=<name>  Organization name (required)
  -name=<name>          Workspace name (required)
  -description=<text>   Workspace description
  -output=<fmt>         Output format: table (default) or json
`

	flags := parseFlagSchemaFromHelp(help)
	if len(flags) != 4 {
		t.Fatalf("expected 4 flags, got %d", len(flags))
	}

	required := map[string]bool{}
	for _, f := range flags {
		required[f.Name] = f.Required
	}

	if !required["organization"] {
		t.Fatal("expected organization to be required")
	}
	if required["description"] {
		t.Fatal("did not expect description to be required")
	}
}

func TestSchemaCommandMissingPath(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SchemaCommand{Meta: newTestMeta(ui)}

	if got := cmd.Run([]string{}); got != 1 {
		t.Fatalf("expected status 1, got %d", got)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "requires a command path") {
		t.Fatalf("expected missing path error, got %q", ui.ErrorWriter.String())
	}
}

func TestSchemaCommandUnknownCommand(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SchemaCommand{Meta: newTestMeta(ui)}

	if got := cmd.Run([]string{"no-such-command"}); got != 1 {
		t.Fatalf("expected status 1, got %d", got)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "unknown command") {
		t.Fatalf("expected unknown command error, got %q", ui.ErrorWriter.String())
	}
}
