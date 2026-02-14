package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// OAuthTokenListCommand is a command to list OAuth tokens
type oauthTokenLister interface {
	List(ctx context.Context, organization string, options *tfe.OAuthTokenListOptions) (*tfe.OAuthTokenList, error)
}

type OAuthTokenListCommand struct {
	Meta
	organization  string
	format        string
	oauthTokenSvc oauthTokenLister
}

// Run executes the oauthtoken list command
func (c *OAuthTokenListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthtoken list")
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

	// List OAuth tokens
	oauthTokens, err := c.oauthTokenService(client).List(client.Context(), c.organization, &tfe.OAuthTokenListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing OAuth tokens: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(oauthTokens.Items) == 0 {
		c.Ui.Output("No OAuth tokens found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Service Provider User", "Has SSH Key", "Created At"}
	var rows [][]string

	for _, ot := range oauthTokens.Items {
		hasSSHKey := "false"
		if ot.HasSSHKey {
			hasSSHKey = "true"
		}

		rows = append(rows, []string{
			ot.ID,
			ot.ServiceProviderUser,
			hasSSHKey,
			ot.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *OAuthTokenListCommand) oauthTokenService(client *client.Client) oauthTokenLister {
	if c.oauthTokenSvc != nil {
		return c.oauthTokenSvc
	}
	return client.OAuthTokens
}

// Help returns help text for the oauthtoken list command
func (c *OAuthTokenListCommand) Help() string {
	helpText := `
Usage: hcptf oauthtoken list [options]

  List OAuth tokens in an organization. OAuth tokens represent
  authenticated connections to a VCS provider for a specific user.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf oauthtoken list -organization=my-org
  hcptf oauthtoken list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthtoken list command
func (c *OAuthTokenListCommand) Synopsis() string {
	return "List OAuth tokens in an organization"
}
