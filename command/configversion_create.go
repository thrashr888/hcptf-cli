package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ConfigVersionCreateCommand is a command to create a configuration version
type ConfigVersionCreateCommand struct {
	Meta
	organization  string
	workspace     string
	autoQueueRuns bool
	speculative   bool
	provisional   bool
	format        string
	workspaceSvc  workspaceReader
	configVerSvc  configVersionCreator
}

// Run executes the configversion create command
func (c *ConfigVersionCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("configversion create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.BoolVar(&c.autoQueueRuns, "auto-queue-runs", true, "Automatically queue runs when uploaded")
	flags.BoolVar(&c.speculative, "speculative", false, "Create a speculative configuration version")
	flags.BoolVar(&c.provisional, "provisional", false, "Create a provisional configuration version")
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

	// Create configuration version
	options := tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: tfe.Bool(c.autoQueueRuns),
		Speculative:   tfe.Bool(c.speculative),
		Provisional:   tfe.Bool(c.provisional),
	}

	configVersion, err := c.configVersionService(client).Create(client.Context(), ws.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating configuration version: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format != "json" {
		c.Ui.Output("Configuration version created successfully")
	}

	data := map[string]interface{}{
		"ID":            configVersion.ID,
		"Status":        configVersion.Status,
		"Source":        configVersion.Source,
		"Speculative":   configVersion.Speculative,
		"Provisional":   configVersion.Provisional,
		"AutoQueueRuns": c.autoQueueRuns,
		"UploadURL":     configVersion.UploadURL,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ConfigVersionCreateCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *ConfigVersionCreateCommand) configVersionService(client *client.Client) configVersionCreator {
	if c.configVerSvc != nil {
		return c.configVerSvc
	}
	return client.ConfigurationVersions
}

// Help returns help text for the configversion create command
func (c *ConfigVersionCreateCommand) Help() string {
	helpText := `
Usage: hcptf configversion create [options]

  Create a new configuration version.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -auto-queue-runs     Automatically queue runs when uploaded (default: true)
  -speculative         Create a speculative configuration version (default: false)
  -provisional         Create a provisional configuration version (default: false)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf configversion create -org=my-org -workspace=my-workspace
  hcptf configversion create -org=my-org -workspace=prod -speculative
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the configversion create command
func (c *ConfigVersionCreateCommand) Synopsis() string {
	return "Create a new configuration version"
}
