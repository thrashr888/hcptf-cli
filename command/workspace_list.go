package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type WorkspaceListCommand struct {
	Meta
	organization string
	format       string
	workspaceSvc workspaceLister
}

// Run executes the workspace list command
func (c *WorkspaceListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List workspaces
	workspaces, err := c.workspaceService(client).List(client.Context(), c.organization, &tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing workspaces: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(workspaces.Items) == 0 {
		c.Ui.Output("No workspaces found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Terraform Version", "Auto Apply", "Locked"}
	var rows [][]string

	for _, ws := range workspaces.Items {
		autoApply := "false"
		if ws.AutoApply {
			autoApply = "true"
		}
		locked := "false"
		if ws.Locked {
			locked = "true"
		}

		rows = append(rows, []string{
			ws.ID,
			ws.Name,
			ws.TerraformVersion,
			autoApply,
			locked,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *WorkspaceListCommand) workspaceService(client *client.Client) workspaceLister {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace list command
func (c *WorkspaceListCommand) Help() string {
	helpText := `
Usage: hcptf workspace list [options]

  List workspaces in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace list -organization=my-org
  hcptf workspace list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace list command
func (c *WorkspaceListCommand) Synopsis() string {
	return "List workspaces in an organization"
}
