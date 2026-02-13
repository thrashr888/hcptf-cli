package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

)

// AccountCreateCommand is a command to create a new user account
type AccountCreateCommand struct {
	Meta
	email    string
	username string
	password string
	format   string
}

// Run executes the account create command
func (c *AccountCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("account create")
	flags.StringVar(&c.email, "email", "", "Email address (required)")
	flags.StringVar(&c.username, "username", "", "Username (required)")
	flags.StringVar(&c.password, "password", "", "Password (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.email == "" {
		c.Ui.Error("Error: -email flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.username == "" {
		c.Ui.Error("Error: -username flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.password == "" {
		c.Ui.Error("Error: -password flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate password strength
	if len(c.password) < 8 {
		c.Ui.Error("Error: password must be at least 8 characters")
		return 1
	}

	// Note: Account creation is an unauthenticated endpoint
	// We'll make a direct API call since go-tfe doesn't expose this yet

	// Get the API address
	address := "https://app.terraform.io"
	if addr := os.Getenv("HCPTF_ADDRESS"); addr != "" {
		address = addr
	}

	// Create account using direct HTTP call
	accountID, accountEmail, accountUsername, err := createAccountDirectly(address, c.email, c.username, c.password)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating account: %s", err))
		c.Ui.Error("")
		c.Ui.Error("Common issues:")
		c.Ui.Error("  - Email already registered")
		c.Ui.Error("  - Username already taken")
		c.Ui.Error("  - Password too weak")
		c.Ui.Error("  - Email verification required")
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output("Account created successfully!")
	c.Ui.Output("")
	c.Ui.Output("Important: Check your email for verification link.")
	c.Ui.Output("")

	// Show account details
	data := map[string]interface{}{
		"ID":       accountID,
		"Email":    accountEmail,
		"Username": accountUsername,
	}

	formatter.KeyValue(data)

	c.Ui.Output("")
	c.Ui.Output("Next steps:")
	c.Ui.Output("  1. Verify your email address")
	c.Ui.Output("  2. Login with: hcptf login")
	c.Ui.Output("  3. Create an organization")

	return 0
}

// Help returns help text for the account create command
func (c *AccountCreateCommand) Help() string {
	helpText := `
Usage: hcptf account create [options]

  Create a new HCP Terraform user account.

  This command creates a new user account on HCP Terraform. After creation,
  you'll need to verify your email address before you can use the account.

Options:

  -email=<email>       Email address (required)
  -username=<username> Username (required)
  -password=<password> Password (required, min 8 characters)
  -output=<format>     Output format: table (default) or json

Example:

  # Create new account
  hcptf account create \
    -email=user@example.com \
    -username=myusername \
    -password=securepassword123

  # After account creation:
  # 1. Check email for verification link
  # 2. Verify email address
  # 3. Login: hcptf login

Note:

  - The email address must not already be registered
  - The username must be unique
  - Password must be at least 8 characters
  - Email verification is required before you can use the account
  - This is for HCP Terraform (app.terraform.io) only
  - Terraform Enterprise accounts are managed by your organization admin
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the account create command
func (c *AccountCreateCommand) Synopsis() string {
	return "Create a new user account"
}

// createAccountDirectly creates an account using direct HTTP API calls
func createAccountDirectly(address, email, username, password string) (string, string, string, error) {
	// Build API request
	requestBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "users",
			"attributes": map[string]string{
				"email":    email,
				"username": username,
				"password": password,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request to account creation endpoint
	url := fmt.Sprintf("%s/api/v2/account/create", address)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", "", "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Data.ID, response.Data.Attributes.Email, response.Data.Attributes.Username, nil
}
