package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// ProjectCreateCommand is a command to create a project
type ProjectCreateCommand struct {
	Meta
	organization string
	name         string
	description  string
	format       string
}

// Run executes the project create command
func (c *ProjectCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("project create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Project name (required)")
	flags.StringVar(&c.description, "description", "", "Project description")
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

	// Create project
	options := tfe.ProjectCreateOptions{
		Name: c.name,
	}

	if c.description != "" {
		options.Description = &c.description
	}

	project, err := client.Projects.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating project: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Project '%s' created successfully", project.Name))

	// Show project details
	data := map[string]interface{}{
		"ID":          project.ID,
		"Name":        project.Name,
		"Description": project.Description,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the project create command
func (c *ProjectCreateCommand) Help() string {
	helpText := `
Usage: hcptf project create [options]

  Create a new project.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Project name (required)
  -description=<text>  Project description
  -output=<format>     Output format: table (default) or json

Example:

  hcptf project create -org=my-org -name=infrastructure
  hcptf project create -org=my-org -name=platform -description="Platform services"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project create command
func (c *ProjectCreateCommand) Synopsis() string {
	return "Create a new project"
}
