package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// WorkspaceTagListCommand is a command to list tags for a workspace
type WorkspaceTagListCommand struct {
	Meta
	workspaceID  string
	organization string
	workspace    string
	format       string
}

// Run executes the workspacetag list command
func (c *WorkspaceTagListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspacetag list")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID")
	flags.StringVar(&c.workspaceID, "id", "", "Workspace ID (alias)")
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

	// List workspace tags
	tags, err := client.Workspaces.ListTags(client.Context(), workspaceID, &tfe.WorkspaceTagListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing workspace tags: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(tags.Items) == 0 {
		c.Ui.Output("No tags found for this workspace")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name"}
	var rows [][]string

	for _, tag := range tags.Items {
		rows = append(rows, []string{
			tag.ID,
			tag.Name,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the workspacetag list command
func (c *WorkspaceTagListCommand) Help() string {
	helpText := `
Usage: hcptf workspacetag list [options]

  List tags applied to a workspace. Tags are organization-level labels
  that help categorize and organize workspaces. Tags must exist in the
  organization before they can be applied to workspaces.

Options:

  -workspace-id=<id>     Workspace ID (format: ws-xxx)
  -id=<id>               Alias for -workspace-id
  -organization=<name>   Organization name
  -org=<name>            Alias for -organization
  -workspace=<name>      Workspace name
  -output=<format>       Output format: table (default) or json

  Either -workspace-id OR both -organization and -workspace are required.

Examples:

  # Using workspace ID
  hcptf workspacetag list -workspace-id=ws-ABC123

  # Using org and workspace name
  hcptf workspacetag list -org=my-org -workspace=prod

  # URL-style
  hcptf my-org prod tags
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspacetag list command
func (c *WorkspaceTagListCommand) Synopsis() string {
	return "List tags for a workspace"
}
