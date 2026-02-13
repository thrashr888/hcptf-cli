package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// StateListCommand is a command to list state versions
type StateListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

// Run executes the state list command
func (c *StateListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("state list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List state versions
	stateVersions, err := client.StateVersions.List(client.Context(), &tfe.StateVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
		Organization: c.organization,
		Workspace:    c.workspace,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing state versions: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(stateVersions.Items) == 0 {
		c.Ui.Output("No state versions found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Serial", "Created At", "Resources"}
	var rows [][]string

	for _, sv := range stateVersions.Items {
		resources := "N/A"
		if sv.ResourcesProcessed {
			resources = fmt.Sprintf("%d", len(sv.Resources))
		}

		rows = append(rows, []string{
			sv.ID,
			fmt.Sprintf("%d", sv.Serial),
			sv.CreatedAt.Format("2006-01-02 15:04:05"),
			resources,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the state list command
func (c *StateListCommand) Help() string {
	helpText := `
Usage: hcptf state list [options]

  List state versions for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf state list -org=my-org -workspace=prod
  hcptf state list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the state list command
func (c *StateListCommand) Synopsis() string {
	return "List state versions for a workspace"
}
