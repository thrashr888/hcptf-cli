package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// WorkspaceCreateCommand is a command to create a workspace
type WorkspaceCreateCommand struct {
	Meta
	organization     string
	name             string
	terraformVersion string
	autoApply        bool
	description      string
	format           string
	workspaceSvc     workspaceCreator
}

// Run executes the workspace create command
func (c *WorkspaceCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
	flags.StringVar(&c.terraformVersion, "terraform-version", "", "Terraform version")
	flags.BoolVar(&c.autoApply, "auto-apply", false, "Enable auto-apply")
	flags.StringVar(&c.description, "description", "", "Workspace description")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create workspace
	options := tfe.WorkspaceCreateOptions{
		Name:      tfe.String(c.name),
		AutoApply: tfe.Bool(c.autoApply),
	}

	if c.terraformVersion != "" {
		options.TerraformVersion = tfe.String(c.terraformVersion)
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	workspace, err := c.workspaceService(client).Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating workspace: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Workspace '%s' created successfully", workspace.Name))
	}

	// Show workspace details
	data := map[string]interface{}{
		"ID":               workspace.ID,
		"Name":             workspace.Name,
		"Organization":     c.organization,
		"TerraformVersion": workspace.TerraformVersion,
		"AutoApply":        workspace.AutoApply,
		"Description":      workspace.Description,
		"CreatedAt":        workspace.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *WorkspaceCreateCommand) workspaceService(client *client.Client) workspaceCreator {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace create command
func (c *WorkspaceCreateCommand) Help() string {
	helpText := `
Usage: hcptf workspace create [options]

  Create a new workspace.

Options:

  -organization=<name>      Organization name (required)
  -org=<name>              Alias for -organization
  -name=<name>             Workspace name (required)
  -terraform-version=<ver> Terraform version to use
  -auto-apply              Enable auto-apply (default: false)
  -description=<text>      Workspace description
  -output=<format>         Output format: table (default) or json

Example:

  hcptf workspace create -org=my-org -name=my-workspace
  hcptf workspace create -org=my-org -name=prod -auto-apply -terraform-version=1.5.0
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace create command
func (c *WorkspaceCreateCommand) Synopsis() string {
	return "Create a new workspace"
}
