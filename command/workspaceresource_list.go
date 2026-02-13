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
	workspaceID  string
	organization string
	workspace    string
	format       string
}

// Run executes the workspace resource list command
func (c *WorkspaceResourceListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspaceresource list")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID")
	flags.StringVar(&c.organization, "organization", "", "Organization name")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags - either workspace-id OR org+workspace
	if c.workspaceID == "" && (c.organization == "" || c.workspace == "") {
		c.Ui.Error("Error: either -workspace-id OR both -organization and -workspace flags are required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// If we don't have workspace ID, look it up from org/workspace name
	workspaceID := c.workspaceID
	if workspaceID == "" {
		ws, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
			return 1
		}
		workspaceID = ws.ID
	}

	// List workspace resources
	resources, err := client.WorkspaceResources.List(client.Context(), workspaceID, &tfe.WorkspaceResourceListOptions{
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

  -workspace-id=<id>     Workspace ID (format: ws-xxx)
  -organization=<name>   Organization name
  -org=<name>            Alias for -organization
  -workspace=<name>      Workspace name
  -output=<format>       Output format: table (default) or json

  Either -workspace-id OR both -organization and -workspace are required.

Examples:

  # Using workspace ID
  hcptf workspaceresource list -workspace-id=ws-abc123

  # Using org and workspace name
  hcptf workspaceresource list -org=my-org -workspace=prod

  # URL-style
  hcptf my-org prod resources
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace resource list command
func (c *WorkspaceResourceListCommand) Synopsis() string {
	return "List resources in workspace state"
}
