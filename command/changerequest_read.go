package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChangeRequestReadCommand shows details of a specific change request
type ChangeRequestReadCommand struct {
	Meta
	id     string
	format string
}

// ChangeRequestReadResponse represents a single change request response
type ChangeRequestReadResponse struct {
	Data ChangeRequest `json:"data"`
}

// Run executes the changerequest read command
func (c *ChangeRequestReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("changerequest read")
	flags.StringVar(&c.id, "id", "", "Change request ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Make direct API call to read change request
	apiURL := fmt.Sprintf("%s/api/v2/change-requests/%s", client.GetAddress(), c.id)

	req, err := http.NewRequestWithContext(client.Context(), "GET", apiURL, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating request: %s", err))
		return 1
	}

	// Get token from client for authorization
	token := client.Token()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error making API request: %s", err))
		return 1
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading response: %s", err))
		return 1
	}

	if resp.StatusCode != http.StatusOK {
		c.Ui.Error(fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body)))
		if resp.StatusCode == http.StatusNotFound {
			c.Ui.Error("\nNote: Change request not found or Change Requests feature is not available.")
			c.Ui.Error("This feature requires HCP Terraform Plus or Enterprise.")
		}
		return 1
	}

	var changeRequest ChangeRequestReadResponse
	if err := json.Unmarshal(body, &changeRequest); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	cr := changeRequest.Data
	status := "Open"
	if cr.Attributes.ArchivedAt != nil {
		status = "Archived"
	}

	data := map[string]interface{}{
		"ID":          cr.ID,
		"Subject":     cr.Attributes.Subject,
		"Message":     cr.Attributes.Message,
		"Status":      status,
		"CreatedAt":   cr.Attributes.CreatedAt,
		"UpdatedAt":   cr.Attributes.UpdatedAt,
		"WorkspaceID": cr.Relationships.Workspace.Data.ID,
	}

	if cr.Attributes.ArchivedBy != nil {
		data["ArchivedBy"] = *cr.Attributes.ArchivedBy
		data["ArchivedAt"] = *cr.Attributes.ArchivedAt
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the changerequest read command
func (c *ChangeRequestReadCommand) Help() string {
	helpText := `
Usage: hcptf changerequest read [options]

  Show details of a specific change request.

  Change requests track workspace to-dos and help teams manage
  compliance and best practices requirements.

  Note: This feature requires HCP Terraform Plus or Enterprise.

Options:

  -id=<id>          Change request ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf changerequest read -id=wscr-abc123
  hcptf changerequest read -id=wscr-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the changerequest read command
func (c *ChangeRequestReadCommand) Synopsis() string {
	return "Show details of a specific change request"
}
