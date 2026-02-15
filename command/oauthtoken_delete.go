package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

type oauthTokenDeleter interface {
	Delete(ctx context.Context, oAuthTokenID string) error
}

// OAuthTokenDeleteCommand is a command to delete an OAuth token.
type OAuthTokenDeleteCommand struct {
	Meta
	id            string
	force         bool
	yes           bool
	oauthTokenSvc oauthTokenDeleter
}

// Run executes the oauthtoken delete command.
func (c *OAuthTokenDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthtoken delete")
	flags.StringVar(&c.id, "id", "", "OAuth token ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")
	flags.BoolVar(&c.yes, "y", false, "Confirm delete without prompt")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	apiClient, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete OAuth token '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.ToLower(strings.TrimSpace(confirmation)) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	if err := c.oauthTokenService(apiClient).Delete(apiClient.Context(), c.id); err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting OAuth token: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("OAuth token '%s' deleted successfully", c.id))
	return 0
}

func (c *OAuthTokenDeleteCommand) oauthTokenService(client *client.Client) oauthTokenDeleter {
	if c.oauthTokenSvc != nil {
		return c.oauthTokenSvc
	}
	return client.OAuthTokens
}

// Help returns help text for the oauthtoken delete command.
func (c *OAuthTokenDeleteCommand) Help() string {
	helpText := `
Usage: hcptf oauthtoken delete [options]

  Delete an OAuth token by ID.

Options:

  -id=<id>          OAuth token ID (required)
  -force            Force delete without confirmation
  -y                Confirm delete without prompt

Example:

  hcptf oauthtoken delete -id=ot-hmAyP66qk2AMVdbJ
  hcptf oauthtoken delete -id=ot-hmAyP66qk2AMVdbJ -force
  hcptf oauthtoken delete -id=ot-hmAyP66qk2AMVdbJ -y
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthtoken delete command.
func (c *OAuthTokenDeleteCommand) Synopsis() string {
	return "Delete an OAuth token"
}
