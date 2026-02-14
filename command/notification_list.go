package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// NotificationListCommand is a command to list notification configurations
type NotificationListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

// Run executes the notification list command
func (c *NotificationListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("notification list")
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

	// Get workspace to obtain its ID
	workspace, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// List notification configurations
	notifications, err := client.NotificationConfigurations.List(client.Context(), workspace.ID, &tfe.NotificationConfigurationListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing notification configurations: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(notifications.Items) == 0 {
		c.Ui.Output("No notification configurations found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Destination Type", "Enabled", "Triggers"}
	var rows [][]string

	for _, nc := range notifications.Items {
		enabled := "false"
		if nc.Enabled {
			enabled = "true"
		}

		triggers := strings.Join(nc.Triggers, ", ")

		rows = append(rows, []string{
			nc.ID,
			nc.Name,
			string(nc.DestinationType),
			enabled,
			triggers,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the notification list command
func (c *NotificationListCommand) Help() string {
	helpText := `
Usage: hcptf notification list [options]

  List notification configurations for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf notification list -org=my-org -workspace=my-workspace
  hcptf notification list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the notification list command
func (c *NotificationListCommand) Synopsis() string {
	return "List notification configurations for a workspace"
}
