package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OrganizationListCommand is a command to list organizations
type OrganizationListCommand struct {
	Meta
	format string
}

// Run executes the organization list command
func (c *OrganizationListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organization list")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List organizations
	orgs, err := client.Organizations.List(client.Context(), &tfe.OrganizationListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing organizations: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(orgs.Items) == 0 {
		c.Ui.Output("No organizations found")
		return 0
	}

	// Prepare table data
	headers := []string{"Name", "Email", "Collaborator Auth Policy", "Created At"}
	var rows [][]string

	for _, org := range orgs.Items {
		rows = append(rows, []string{
			org.Name,
			org.Email,
			string(org.CollaboratorAuthPolicy),
			org.CreatedAt.Format("2006-01-02"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the organization list command
func (c *OrganizationListCommand) Help() string {
	helpText := `
Usage: hcptf organization list [options]

  List organizations accessible to the authenticated user.

Options:

  -output=<format>  Output format: table (default) or json

Example:

  hcptf organization list
  hcptf organization list -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization list command
func (c *OrganizationListCommand) Synopsis() string {
	return "List organizations"
}
