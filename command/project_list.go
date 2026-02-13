package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ProjectListCommand is a command to list projects
type ProjectListCommand struct {
	Meta
	organization string
	format       string
	projectSvc   projectLister
}

// Run executes the project list command
func (c *ProjectListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("project list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List projects
	projects, err := c.projectService(client).List(client.Context(), c.organization, &tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing projects: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(projects.Items) == 0 {
		c.Ui.Output("No projects found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Description"}
	var rows [][]string

	for _, project := range projects.Items {
		description := project.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}

		rows = append(rows, []string{
			project.ID,
			project.Name,
			description,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the project list command
func (c *ProjectListCommand) Help() string {
	helpText := `
Usage: hcptf project list [options]

  List projects in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf project list -org=my-org
  hcptf project list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *ProjectListCommand) projectService(client *client.Client) projectLister {
	if c.projectSvc != nil {
		return c.projectSvc
	}
	return client.Projects
}

// Synopsis returns a short synopsis for the project list command
func (c *ProjectListCommand) Synopsis() string {
	return "List projects in an organization"
}
