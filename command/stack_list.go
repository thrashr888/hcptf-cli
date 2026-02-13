package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// StackListCommand is a command to list stacks
type StackListCommand struct {
	Meta
	organization string
	project      string
	format       string
}

// Run executes the stack list command
func (c *StackListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stack list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.project, "project", "", "Filter by project ID")
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

	// List stacks
	options := &tfe.StackListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	}

	if c.project != "" {
		options.ProjectID = c.project
	}

	stacks, err := client.Stacks.List(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing stacks: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(stacks.Items) == 0 {
		c.Ui.Output("No stacks found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Description", "Project", "Created"}
	var rows [][]string

	for _, stack := range stacks.Items {
		description := stack.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}

		projectID := ""
		if stack.Project != nil {
			projectID = stack.Project.ID
		}

		rows = append(rows, []string{
			stack.ID,
			stack.Name,
			description,
			projectID,
			stack.CreatedAt.Format("2006-01-02"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the stack list command
func (c *StackListCommand) Help() string {
	helpText := `
Usage: hcptf stack list [options]

  List stacks in an organization or project. Stacks enable orchestrating
  deployments across multiple configurations and workspaces.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -project=<id>        Filter by project ID
  -output=<format>     Output format: table (default) or json

Example:

  hcptf stack list -org=my-org
  hcptf stack list -org=my-org -project=prj-abc123
  hcptf stack list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack list command
func (c *StackListCommand) Synopsis() string {
	return "List stacks in an organization or project"
}
