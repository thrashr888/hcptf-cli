package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

)

// ChangeRequestUpdateCommand archives a change request
type ChangeRequestUpdateCommand struct {
	Meta
	id      string
	archive bool
	format  string
}

// Run executes the changerequest update command
func (c *ChangeRequestUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("changerequest update")
	flags.StringVar(&c.id, "id", "", "Change request ID (required)")
	flags.BoolVar(&c.archive, "archive", false, "Archive the change request")
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

	if !c.archive {
		c.Ui.Error("Error: -archive flag is required (currently only archive operation is supported)")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Make direct API call to archive change request
	apiURL := fmt.Sprintf("%s/api/v2/change-requests/%s", client.GetAddress(), c.id)

	req, err := http.NewRequestWithContext(client.Context(), "PATCH", apiURL, nil)
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

	c.Ui.Output(fmt.Sprintf("Change request '%s' archived successfully", c.id))

	cr := changeRequest.Data
	data := map[string]interface{}{
		"ID":        cr.ID,
		"Subject":   cr.Attributes.Subject,
		"Status":    "Archived",
		"CreatedAt": cr.Attributes.CreatedAt,
		"UpdatedAt": cr.Attributes.UpdatedAt,
	}

	if cr.Attributes.ArchivedBy != nil {
		data["ArchivedBy"] = *cr.Attributes.ArchivedBy
		data["ArchivedAt"] = *cr.Attributes.ArchivedAt
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the changerequest update command
func (c *ChangeRequestUpdateCommand) Help() string {
	helpText := `
Usage: hcptf changerequest update [options]

  Update a change request (currently only archiving is supported).

  Archiving a change request marks it as complete. Archived change
  requests remain visible but are sorted separately from active requests.

  Note: This feature requires HCP Terraform Plus or Enterprise.

Options:

  -id=<id>          Change request ID (required)
  -archive          Archive the change request (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf changerequest update -id=wscr-abc123 -archive
  hcptf changerequest update -id=wscr-abc123 -archive -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the changerequest update command
func (c *ChangeRequestUpdateCommand) Synopsis() string {
	return "Update a change request (archive)"
}
