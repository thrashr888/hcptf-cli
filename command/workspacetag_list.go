package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// WorkspaceTagListCommand is a command to list tags for a workspace
type WorkspaceTagListCommand struct {
	Meta
	workspaceID string
	format      string
}

// Run executes the workspacetag list command
func (c *WorkspaceTagListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspacetag list")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.workspaceID, "id", "", "Workspace ID (alias)")
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

	// List workspace tags
	tags, err := client.Workspaces.ListTags(client.Context(), c.workspaceID, &tfe.WorkspaceTagListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing workspace tags: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

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

  -workspace-id=<id>   Workspace ID (required)
  -id=<id>            Alias for -workspace-id
  -output=<format>    Output format: table (default) or json

Example:

  hcptf workspacetag list -workspace-id=ws-ABC123
  hcptf workspacetag list -id=ws-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspacetag list command
func (c *WorkspaceTagListCommand) Synopsis() string {
	return "List tags for a workspace"
}
