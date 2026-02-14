package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// ProjectUpdateCommand is a command to update a project
type ProjectUpdateCommand struct {
	Meta
	projectID   string
	name        string
	description string
	format      string
}

// Run executes the project update command
func (c *ProjectUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("project update")
	flags.StringVar(&c.projectID, "id", "", "Project ID (required)")
	flags.StringVar(&c.name, "name", "", "New project name")
	flags.StringVar(&c.description, "description", "", "New project description")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

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

	// Build update options
	options := tfe.ProjectUpdateOptions{}

	if c.name != "" {
		options.Name = &c.name
	}

	if c.description != "" {
		options.Description = &c.description
	}

	// Update project
	project, err := client.Projects.Update(client.Context(), c.projectID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating project: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Project '%s' updated successfully", project.Name))

	// Show project details
	data := map[string]interface{}{
		"ID":          project.ID,
		"Name":        project.Name,
		"Description": project.Description,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the project update command
func (c *ProjectUpdateCommand) Help() string {
	helpText := `
Usage: hcptf project update [options]

  Update project settings.

Options:

  -id=<project-id>     Project ID (required)
  -name=<name>         New project name
  -description=<text>  New project description
  -output=<format>     Output format: table (default) or json

Example:

  hcptf project update -id=prj-abc123 -name="New Name"
  hcptf project update -id=prj-abc123 -description="Updated description"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project update command
func (c *ProjectUpdateCommand) Synopsis() string {
	return "Update project settings"
}
