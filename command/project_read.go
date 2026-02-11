package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ProjectReadCommand is a command to read project details
type ProjectReadCommand struct {
	Meta
	projectID string
	format    string
}

// Run executes the project read command
func (c *ProjectReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("project read")
	flags.StringVar(&c.projectID, "id", "", "Project ID (required)")
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

	// Read project
	project, err := client.Projects.Read(client.Context(), c.projectID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading project: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":          project.ID,
		"Name":        project.Name,
		"Description": project.Description,
	}

	if project.Organization != nil {
		data["Organization"] = project.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the project read command
func (c *ProjectReadCommand) Help() string {
	helpText := `
Usage: hcptf project read [options]

  Read project details.

Options:

  -id=<project-id>  Project ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf project read -id=prj-abc123
  hcptf project read -id=prj-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project read command
func (c *ProjectReadCommand) Synopsis() string {
	return "Read project details"
}
