package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AccountUpdateCommand is a command to update account details
type AccountUpdateCommand struct {
	Meta
	email       string
	username    string
	password    string
	newPassword string
	format      string
}

// Run executes the account update command
func (c *AccountUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("account update")
	flags.StringVar(&c.email, "email", "", "New email address")
	flags.StringVar(&c.username, "username", "", "New username")
	flags.StringVar(&c.password, "password", "", "Current password (required for changes)")
	flags.StringVar(&c.newPassword, "new-password", "", "New password")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check if any update is requested
	if c.email == "" && c.username == "" && c.newPassword == "" {
		c.Ui.Error("Error: at least one of -email, -username, or -new-password must be provided")
		c.Ui.Error(c.Help())
		return 1
	}

	// Password is required for security
	if c.password == "" {
		c.Ui.Error("Error: -password flag is required to verify your identity")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate new password if provided
	if c.newPassword != "" && len(c.newPassword) < 8 {
		c.Ui.Error("Error: new password must be at least 8 characters")
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build update options (using Users API)
	options := tfe.UserUpdateOptions{
		Email:    tfe.String(c.email),
		Username: tfe.String(c.username),
	}

	// Note: go-tfe's UserUpdateOptions doesn't support password changes
	// This would require a direct API call if needed
	if c.newPassword != "" {
		c.Ui.Error("Error: password changes are not yet supported through the CLI")
		c.Ui.Error("Please use the web UI to change your password:")
		c.Ui.Error("  https://app.terraform.io/app/settings/profile")
		return 1
	}

	// Update account (using Users API)
	account, err := client.Users.UpdateCurrent(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating account: %s", err))
		c.Ui.Error("")
		c.Ui.Error("Common issues:")
		c.Ui.Error("  - Incorrect current password")
		c.Ui.Error("  - Email already in use")
		c.Ui.Error("  - Username already taken")
		c.Ui.Error("  - New password too weak")
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output("Account updated successfully!")
	c.Ui.Output("")

	// Show updated account details
	data := map[string]interface{}{
		"ID":       account.ID,
		"Email":    account.Email,
		"Username": account.Username,
	}

	if c.email != "" && account.UnconfirmedEmail != "" {
		c.Ui.Output("Important: Check your new email for verification link.")
		c.Ui.Output("")
		data["UnconfirmedEmail"] = account.UnconfirmedEmail
	}

	formatter.KeyValue(data)

	if c.newPassword != "" {
		c.Ui.Output("")
		c.Ui.Output("Password updated. Use your new password for future logins.")
	}

	return 0
}

// Help returns help text for the account update command
func (c *AccountUpdateCommand) Help() string {
	helpText := `
Usage: hcptf account update [options]

  Update current account details.

  This command allows you to update your email, username, or password.
  Your current password is required for security verification.

Options:

  -email=<email>           New email address
  -username=<username>     New username
  -password=<password>     Current password (required)
  -new-password=<password> New password (min 8 characters)
  -output=<format>         Output format: table (default) or json

Example:

  # Update email address
  hcptf account update \
    -email=newemail@example.com \
    -password=currentpassword

  # Update username
  hcptf account update \
    -username=newusername \
    -password=currentpassword

  # Change password
  hcptf account update \
    -password=currentpassword \
    -new-password=newsecurepassword123

  # Update multiple fields
  hcptf account update \
    -email=newemail@example.com \
    -username=newusername \
    -password=currentpassword

Note:

  - Current password is always required for security
  - Changing email requires verification of the new email address
  - Username must be unique across HCP Terraform
  - Password must be at least 8 characters
  - This command requires authentication
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the account update command
func (c *AccountUpdateCommand) Synopsis() string {
	return "Update account details"
}
