package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// WorkspaceResourceListCommand is a command to list workspace resources
type WorkspaceResourceListCommand struct {
	Meta
	workspaceID string
	format      string
}

// Run executes the workspace resource list command
func (c *WorkspaceResourceListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspaceresource list")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.workspaceID == "" {
		c.Ui.Error("Error: -workspace-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List workspace resources
	resources, err := client.WorkspaceResources.List(client.Context(), c.workspaceID, &tfe.WorkspaceResourceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing workspace resources: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(resources.Items) == 0 {
		c.Ui.Output("No resources found in workspace state")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Address", "Type", "Provider", "Module"}
	var rows [][]string

	// Count resources by type for summary
	typeCounts := make(map[string]int)

	for _, resource := range resources.Items {
		rows = append(rows, []string{
			resource.ID,
			resource.Address,
			resource.ProviderType,
			resource.Provider,
			resource.Module,
		})
		typeCounts[resource.ProviderType]++
	}

	// Display summary if table format
	if c.format == "table" {
		c.Ui.Output(fmt.Sprintf("Total resources: %d\n", len(resources.Items)))

		if len(typeCounts) > 0 {
			c.Ui.Output("Resource types:")
			for resType, count := range typeCounts {
				c.Ui.Output(fmt.Sprintf("  %s: %d", resType, count))
			}
			c.Ui.Output("") // blank line before table
		}
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the workspace resource list command
func (c *WorkspaceResourceListCommand) Help() string {
	helpText := `
Usage: hcptf workspaceresource list [options]

  List resources in workspace state. Shows all infrastructure resources
  currently managed by the workspace.

Options:

  -workspace-id=<id>   Workspace ID (required, format: ws-xxx)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspaceresource list -workspace-id=ws-abc123
  hcptf workspaceresource list -workspace-id=ws-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace resource list command
func (c *WorkspaceResourceListCommand) Synopsis() string {
	return "List resources in workspace state"
}
