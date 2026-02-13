package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// QueryWorkspaceListCommand is a command to search workspaces across an organization
type QueryWorkspaceListCommand struct {
	Meta
	organization string
	search       string
	tags         string
	excludeTags  string
	wildcard     string
	format       string
}

// Run executes the query workspace list command
func (c *QueryWorkspaceListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("queryworkspace list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.search, "search", "", "Search query for workspace name")
	flags.StringVar(&c.tags, "tags", "", "Filter by tags (comma-separated)")
	flags.StringVar(&c.excludeTags, "exclude-tags", "", "Exclude workspaces with tags (comma-separated)")
	flags.StringVar(&c.wildcard, "wildcard", "", "Wildcard filter for workspace name")
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

	// Build list options
	options := &tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
		Include: []tfe.WSIncludeOpt{
			tfe.WSCurrentRun,
			tfe.WSOrganization,
		},
	}

	if c.search != "" {
		options.Search = c.search
	}

	if c.tags != "" {
		options.Tags = c.tags
	}

	if c.excludeTags != "" {
		options.ExcludeTags = c.excludeTags
	}

	if c.wildcard != "" {
		options.WildcardName = c.wildcard
	}

	// List workspaces
	workspaces, err := client.Workspaces.List(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing workspaces: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(workspaces.Items) == 0 {
		c.Ui.Output("No workspaces found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Terraform Version", "Current Run", "Auto Apply", "Locked"}
	var rows [][]string

	for _, ws := range workspaces.Items {
		currentRunStatus := "None"
		if ws.CurrentRun != nil {
			currentRunStatus = string(ws.CurrentRun.Status)
		}

		rows = append(rows, []string{
			ws.ID,
			ws.Name,
			ws.TerraformVersion,
			currentRunStatus,
			fmt.Sprintf("%v", ws.AutoApply),
			fmt.Sprintf("%v", ws.Locked),
		})
	}

	formatter.Table(headers, rows)

	// Show summary
	c.Ui.Output(fmt.Sprintf("\nTotal: %d workspace(s)", len(workspaces.Items)))

	return 0
}

// Help returns help text for the query workspace list command
func (c *QueryWorkspaceListCommand) Help() string {
	helpText := `
Usage: hcptf queryworkspace list [options]

  Search and list workspaces across an organization with filters.

  This command allows you to search for workspaces in an organization using
  various filters including name search, tags, and wildcard patterns.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -search=<query>      Search query for workspace name
  -tags=<tags>         Filter by tags (comma-separated, e.g., "env:prod,team:platform")
  -exclude-tags=<tags> Exclude workspaces with these tags (comma-separated)
  -wildcard=<pattern>  Wildcard filter for workspace name (e.g., "prod-*")
  -output=<format>     Output format: table (default) or json

Example:

  # List all workspaces in organization
  hcptf queryworkspace list -org=my-org

  # Search for workspaces by name
  hcptf queryworkspace list -org=my-org -search=production

  # Filter by tags
  hcptf queryworkspace list -org=my-org -tags=env:prod

  # Filter by tags and exclude others
  hcptf queryworkspace list -org=my-org -tags=env:prod -exclude-tags=archived

  # Use wildcard pattern
  hcptf queryworkspace list -org=my-org -wildcard="prod-*"

  # Output as JSON
  hcptf queryworkspace list -org=my-org -output=json

Note:

  Workspaces are returned with information about their current run status,
  Terraform version, and configuration settings.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the query workspace list command
func (c *QueryWorkspaceListCommand) Synopsis() string {
	return "Search workspaces across organization"
}
