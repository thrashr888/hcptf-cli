package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationTagListCommand is a command to list organization tags
type OrganizationTagListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the organizationtag list command
func (c *OrganizationTagListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtag list")
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

	// List organization tags
	tags, err := client.OrganizationTags.List(client.Context(), c.organization, &tfe.OrganizationTagsListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing organization tags: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(tags.Items) == 0 {
		c.Ui.Output("No organization tags found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Instance Count"}
	var rows [][]string

	for _, tag := range tags.Items {
		rows = append(rows, []string{
			tag.ID,
			tag.Name,
			fmt.Sprintf("%d", tag.InstanceCount),
			"-", // CreatedAt not available in go-tfe
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the organizationtag list command
func (c *OrganizationTagListCommand) Help() string {
	helpText := `
Usage: hcptf organizationtag list [options]

  List tags used in workspaces across the organization.
  Organization tags are key-value pairs for categorizing resources
  like workspaces. Tags can be used for filtering and organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf organizationtag list -org=my-org
  hcptf organizationtag list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organizationtag list command
func (c *OrganizationTagListCommand) Synopsis() string {
	return "List organization tags"
}
