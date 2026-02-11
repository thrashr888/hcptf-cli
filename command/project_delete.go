package command

import (
	"fmt"
	"strings"
)

// ProjectDeleteCommand is a command to delete a project
type ProjectDeleteCommand struct {
	Meta
	projectID string
	force     bool
}

// Run executes the project delete command
func (c *ProjectDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("project delete")
	flags.StringVar(&c.projectID, "id", "", "Project ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.projectID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete project '%s'? (yes/no): ", c.projectID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete project
	err = client.Projects.Delete(client.Context(), c.projectID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting project: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Project '%s' deleted successfully", c.projectID))
	return 0
}

// Help returns help text for the project delete command
func (c *ProjectDeleteCommand) Help() string {
	helpText := `
Usage: hcptf project delete [options]

  Delete a project.

Options:

  -id=<project-id>  Project ID (required)
  -force            Force delete without confirmation

Example:

  hcptf project delete -id=prj-abc123
  hcptf project delete -id=prj-old123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project delete command
func (c *ProjectDeleteCommand) Synopsis() string {
	return "Delete a project"
}
