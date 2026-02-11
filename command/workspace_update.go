package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// WorkspaceUpdateCommand is a command to update a workspace
type WorkspaceUpdateCommand struct {
	Meta
	organization     string
	name             string
	newName          string
	terraformVersion string
	autoApply        string
	description      string
	format           string
	workspaceSvc     workspaceUpdater
}

// Run executes the workspace update command
func (c *WorkspaceUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace update")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
	flags.StringVar(&c.newName, "new-name", "", "New workspace name")
	flags.StringVar(&c.terraformVersion, "terraform-version", "", "Terraform version")
	flags.StringVar(&c.autoApply, "auto-apply", "", "Enable auto-apply (true/false)")
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

	// Build update options
	options := tfe.WorkspaceUpdateOptions{}

	if c.newName != "" {
		options.Name = tfe.String(c.newName)
	}

	if c.terraformVersion != "" {
		options.TerraformVersion = tfe.String(c.terraformVersion)
	}

	if c.autoApply != "" {
		if c.autoApply == "true" {
			options.AutoApply = tfe.Bool(true)
		} else if c.autoApply == "false" {
			options.AutoApply = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -auto-apply must be 'true' or 'false'")
			return 1
		}
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	// Update workspace
	workspace, err := c.workspaceService(client).Update(client.Context(), c.organization, c.name, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating workspace: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Workspace '%s' updated successfully", workspace.Name))

	// Show workspace details
	data := map[string]interface{}{
		"ID":               workspace.ID,
		"Name":             workspace.Name,
		"Organization":     c.organization,
		"TerraformVersion": workspace.TerraformVersion,
		"AutoApply":        workspace.AutoApply,
		"Description":      workspace.Description,
		"UpdatedAt":        workspace.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *WorkspaceUpdateCommand) workspaceService(client *client.Client) workspaceUpdater {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace update command
func (c *WorkspaceUpdateCommand) Help() string {
	helpText := `
Usage: hcptf workspace update [options]

  Update workspace settings.

Options:

  -organization=<name>      Organization name (required)
  -org=<name>              Alias for -organization
  -name=<name>             Workspace name (required)
  -new-name=<name>         New workspace name
  -terraform-version=<ver> Terraform version to use
  -auto-apply=<bool>       Enable auto-apply (true/false)
  -description=<text>      Workspace description
  -output=<format>         Output format: table (default) or json

Example:

  hcptf workspace update -org=my-org -name=my-workspace -auto-apply=true
  hcptf workspace update -org=my-org -name=old-name -new-name=new-name
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace update command
func (c *WorkspaceUpdateCommand) Synopsis() string {
	return "Update workspace settings"
}
