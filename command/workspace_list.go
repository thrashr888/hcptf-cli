package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type WorkspaceListCommand struct {
	Meta
	organization     string
	search           string
	tags             string
	excludeTags      string
	wildcardName     string
	projectID        string
	currentRunStatus string
	include          string
	sort             string
	format           string
	workspaceSvc     workspaceLister
}

// Run executes the workspace list command
func (c *WorkspaceListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.search, "search", "", "Search workspace names by substring")
	flags.StringVar(&c.tags, "tags", "", "Filter by tag names (comma-separated)")
	flags.StringVar(&c.excludeTags, "exclude-tags", "", "Exclude workspaces with these tags (comma-separated)")
	flags.StringVar(&c.wildcardName, "wildcard-name", "", "Wildcard filter for workspace names")
	flags.StringVar(&c.projectID, "project-id", "", "Filter by project ID")
	flags.StringVar(&c.currentRunStatus, "current-run-status", "", "Filter by current run status")
	flags.StringVar(&c.include, "include", "", "Comma-separated related resources to include")
	flags.StringVar(&c.sort, "sort", "", "Sort order (e.g. name,-name,current-run.created-at)")
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

	options := &tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
		Search:           c.search,
		Tags:             c.tags,
		ExcludeTags:      c.excludeTags,
		WildcardName:     c.wildcardName,
		ProjectID:        c.projectID,
		CurrentRunStatus: c.currentRunStatus,
		Sort:             c.sort,
	}
	if c.include != "" {
		for _, include := range splitCommaList(c.include) {
			if include == "" {
				continue
			}
			options.Include = append(options.Include, tfe.WSIncludeOpt(include))
		}
	}

	// List workspaces
	workspaces, err := c.workspaceService(client).List(client.Context(), c.organization, options)
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
	headers := []string{"ID", "Name", "Terraform Version", "Auto Apply", "Locked"}
	var rows [][]string

	for _, ws := range workspaces.Items {
		autoApply := "false"
		if ws.AutoApply {
			autoApply = "true"
		}
		locked := "false"
		if ws.Locked {
			locked = "true"
		}

		rows = append(rows, []string{
			ws.ID,
			ws.Name,
			ws.TerraformVersion,
			autoApply,
			locked,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *WorkspaceListCommand) workspaceService(client *client.Client) workspaceLister {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace list command
func (c *WorkspaceListCommand) Help() string {
	helpText := `
Usage: hcptf workspace list [options]

  List workspaces in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -search=<query>      Search workspace names by substring
  -tags=<names>        Filter by tag names (comma-separated)
  -exclude-tags=<names> Exclude workspaces with tags (comma-separated)
  -wildcard-name=<pat> Wildcard filter for workspace names
  -project-id=<id>     Filter by project ID
  -current-run-status=<status> Filter by current run status
  -include=<values>    Comma-separated include values
  -sort=<value>        Sort order (e.g. name,-name,current-run.created-at)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace list -organization=my-org
  hcptf workspace list -org=my-org -search=prod -tags=env:prod
  hcptf workspace list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace list command
func (c *WorkspaceListCommand) Synopsis() string {
	return "List workspaces in an organization"
}
