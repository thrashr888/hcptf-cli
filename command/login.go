package command

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/config"
)

// LoginCommand is a command to authenticate and save credentials
type LoginCommand struct {
	Meta
	hostname  string
	showToken bool
}

// Run executes the login command
func (c *LoginCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("login")
	flags.StringVar(&c.hostname, "hostname", "app.terraform.io", "HCP Terraform hostname")
	flags.BoolVar(&c.showToken, "show-token", false, "Show token after successful login")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.showToken {
		cfg, err := c.Config()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error loading config: %s", err))
			return 1
		}

		token := cfg.GetToken(c.hostname)
		if token == "" {
			c.Ui.Error("No token found for this hostname. Run 'hcptf login' first.")
			return 1
		}

		c.Ui.Output(token)
		return 0
	}

	credsPath := config.GetTerraformCredentialsPath()

	c.Ui.Output(fmt.Sprintf("hcptf will request an API token for %s using your browser.", c.hostname))
	c.Ui.Output("")
	c.Ui.Output("If login is successful, hcptf will store the token in plain text in")
	c.Ui.Output("the following file for use by subsequent commands:")
	c.Ui.Output(fmt.Sprintf("    %s", credsPath))
	c.Ui.Output("")

	confirm, err := c.Ui.Ask("Do you want to proceed?\n  Only 'yes' will be accepted to confirm.\n\n  Enter a value:")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading input: %s", err))
		return 1
	}
	if strings.TrimSpace(confirm) != "yes" {
		c.Ui.Output("")
		c.Ui.Output("Login cancelled.")
		return 0
	}

	c.Ui.Output("")
	c.Ui.Output(strings.Repeat("-", 81))
	c.Ui.Output("")

	tokensURL := fmt.Sprintf("https://%s/app/settings/tokens?source=terraform-login", c.hostname)

	c.Ui.Output(fmt.Sprintf("hcptf must now open a web browser to the tokens page for %s.", c.hostname))
	c.Ui.Output("")
	c.Ui.Output("If a browser does not open this automatically, open the following URL to proceed:")
	c.Ui.Output(fmt.Sprintf("    %s", tokensURL))

	// Attempt to open browser (best-effort)
	openBrowser(tokensURL)

	c.Ui.Output("")
	c.Ui.Output(strings.Repeat("-", 81))
	c.Ui.Output("")
	c.Ui.Output("Generate a token using your browser, and copy-paste it into this prompt.")
	c.Ui.Output("")
	c.Ui.Output("hcptf will store the token in plain text in the following file")
	c.Ui.Output("for use by subsequent commands:")
	c.Ui.Output(fmt.Sprintf("    %s", credsPath))
	c.Ui.Output("")

	token, err := c.Ui.Ask(fmt.Sprintf("Token for %s:\n  Enter a value:", c.hostname))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading token: %s", err))
		return 1
	}

	token = strings.TrimSpace(token)
	if token == "" {
		c.Ui.Error("Error: token cannot be empty")
		return 1
	}

	// Validate token and retrieve the username
	username, err := config.ValidateTokenAndGetUser(c.hostname, token)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: token validation failed: %s", err))
		c.Ui.Error("")
		c.Ui.Error("Please verify:")
		c.Ui.Error("  1. The token is correct")
		c.Ui.Error("  2. The token has not expired")
		c.Ui.Error(fmt.Sprintf("  3. You can access %s", c.hostname))
		return 1
	}

	// Save token to credentials file
	if err := config.SaveCredential(c.hostname, token); err != nil {
		c.Ui.Error(fmt.Sprintf("Error saving credentials: %s", err))
		return 1
	}

	c.Ui.Output("")
	c.Ui.Output(fmt.Sprintf("Retrieved token for user %s", username))
	c.Ui.Output("")
	c.Ui.Output(strings.Repeat("-", 81))
	c.Ui.Output("")
	c.Ui.Output(hcpTerraformBanner)
	c.Ui.Output("")

	return 0
}

// openBrowser attempts to open the given URL in the default browser.
// It is best-effort and silently ignores errors.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return
	}
	_ = cmd.Start()
}

const hcpTerraformBanner = `                                          -
                                          -----                           -
                                          ---------                      --
                                          ---------  -                -----
                                           ---------  ------        -------
                                             -------  ---------  ----------
                                                ----  ---------- ----------
                                                  --  ---------- ----------
   Welcome to HCP Terraform!                       -  ---------- -------
                                                      ---  ----- ---
   Documentation: terraform.io/docs/cloud             --------   -
                                                      ----------
                                                      ----------
                                                       ---------
                                                           -----
                                                               -`

// Help returns help text for the login command
func (c *LoginCommand) Help() string {
	helpText := `
Usage: hcptf login [options]

  Authenticate to HCP Terraform and save credentials.

  This command will open your browser to generate an API token, then store
  it in ~/.terraform.d/credentials.tfrc.json, making it available to both
  hcptf and the Terraform CLI.

Options:

  -hostname=<hostname>  HCP Terraform hostname (default: app.terraform.io)
  -show-token          Show token and exit without prompting (default: false)

Example:

  # Login to HCP Terraform
  hcptf login

  # Login to Terraform Enterprise
  hcptf login -hostname=tfe.example.com

Note:

  The stored credentials are compatible with 'terraform login' and will
  be automatically used by both the Terraform CLI and hcptf.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the login command
func (c *LoginCommand) Synopsis() string {
	return "Authenticate to HCP Terraform"
}
