package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OAuthClientListCommand is a command to list OAuth clients
type oauthClientLister interface {
	List(ctx context.Context, organization string, options *tfe.OAuthClientListOptions) (*tfe.OAuthClientList, error)
}

type OAuthClientListCommand struct {
	Meta
	organization   string
	format         string
	oauthClientSvc oauthClientLister
}

// Run executes the oauthclient list command
func (c *OAuthClientListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthclient list")
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

	// List OAuth clients
	oauthClients, err := c.oauthClientService(client).List(client.Context(), c.organization, &tfe.OAuthClientListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing OAuth clients: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(oauthClients.Items) == 0 {
		c.Ui.Output("No OAuth clients found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Service Provider", "HTTP URL", "Created At"}
	var rows [][]string

	for _, oc := range oauthClients.Items {
		name := ""
		if oc.Name != nil && *oc.Name != "" {
			name = *oc.Name
		} else {
			name = oc.ServiceProviderName
		}

		rows = append(rows, []string{
			oc.ID,
			name,
			string(oc.ServiceProvider),
			oc.HTTPURL,
			oc.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *OAuthClientListCommand) oauthClientService(client *client.Client) oauthClientLister {
	if c.oauthClientSvc != nil {
		return c.oauthClientSvc
	}
	return client.OAuthClients
}

// Help returns help text for the oauthclient list command
func (c *OAuthClientListCommand) Help() string {
	helpText := `
Usage: hcptf oauthclient list [options]

  List OAuth clients in an organization. OAuth clients represent
  VCS connections (GitHub, GitLab, etc.) used by workspaces.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf oauthclient list -organization=my-org
  hcptf oauthclient list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthclient list command
func (c *OAuthClientListCommand) Synopsis() string {
	return "List OAuth clients in an organization"
}
