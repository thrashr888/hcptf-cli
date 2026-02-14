package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ConfigVersionListCommand is a command to list configuration versions
type ConfigVersionListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
	workspaceSvc workspaceReader
	configVerSvc configVersionLister
}

// Run executes the configversion list command
func (c *ConfigVersionListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("configversion list")
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

	// Get workspace first
	ws, err := c.workspaceService(client).Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// List configuration versions
	configVersions, err := c.configVersionService(client).List(client.Context(), ws.ID, &tfe.ConfigurationVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing configuration versions: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(configVersions.Items) == 0 {
		c.Ui.Output("No configuration versions found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Status", "Source", "Speculative", "Provisional"}
	var rows [][]string

	for _, cv := range configVersions.Items {
		source := string(cv.Source)
		rows = append(rows, []string{
			cv.ID,
			string(cv.Status),
			source,
			fmt.Sprintf("%t", cv.Speculative),
			fmt.Sprintf("%t", cv.Provisional),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *ConfigVersionListCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *ConfigVersionListCommand) configVersionService(client *client.Client) configVersionLister {
	if c.configVerSvc != nil {
		return c.configVerSvc
	}
	return client.ConfigurationVersions
}

// Help returns help text for the configversion list command
func (c *ConfigVersionListCommand) Help() string {
	helpText := `
Usage: hcptf configversion list [options]

  List configuration versions for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf configversion list -org=my-org -workspace=my-workspace
  hcptf configversion list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the configversion list command
func (c *ConfigVersionListCommand) Synopsis() string {
	return "List configuration versions for a workspace"
}
