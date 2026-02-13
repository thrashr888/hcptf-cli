package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

)

// VCSEventReadCommand shows details of a specific VCS event
type VCSEventReadCommand struct {
	Meta
	id     string
	format string
}

// VCSEventReadResponse represents a single VCS event response
type VCSEventReadResponse struct {
	Data VCSEvent `json:"data"`
}

// Run executes the vcsevent read command
func (c *VCSEventReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vcsevent read")
	flags.StringVar(&c.id, "id", "", "VCS Event ID (required)")
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

	// Make direct API call to read VCS event
	apiURL := fmt.Sprintf("%s/api/v2/vcs-events/%s", client.GetAddress(), c.id)

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
			c.Ui.Error("\nNote: VCS event not found or VCS Events feature is not available.")
			c.Ui.Error("The VCS Events API is in beta and currently only supports GitLab.com")
			c.Ui.Error("connections established after December 2020.")
		}
		return 1
	}

	var vcsEvent VCSEventReadResponse
	if err := json.Unmarshal(body, &vcsEvent); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	event := vcsEvent.Data
	data := map[string]interface{}{
		"ID":             event.ID,
		"Level":          event.Attributes.Level,
		"Message":        event.Attributes.Message,
		"OrganizationID": event.Attributes.OrganizationID,
		"CreatedAt":      event.Attributes.CreatedAt,
	}

	// Add optional fields if present
	if event.Attributes.SuggestedAction != nil {
		data["SuggestedAction"] = *event.Attributes.SuggestedAction
	}

	if event.Relationships.OAuthClient != nil && event.Relationships.OAuthClient.Data != nil {
		data["OAuthClientID"] = event.Relationships.OAuthClient.Data.ID
	}

	if event.Relationships.OAuthToken != nil && event.Relationships.OAuthToken.Data != nil {
		data["OAuthTokenID"] = event.Relationships.OAuthToken.Data.ID
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the vcsevent read command
func (c *VCSEventReadCommand) Help() string {
	helpText := `
Usage: hcptf vcsevent read [options]

  Show details of a specific VCS event.

  VCS events describe changes and actions related to VCS integration,
  including webhook deliveries, OAuth token issues, and connection
  problems. This command displays the full event details including
  any suggested actions to resolve issues.

  Note: The VCS Events API is in beta and currently only supports
  GitLab.com connections established after December 2020.

Options:

  -id=<id>          VCS Event ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf vcsevent read -id=ve-DJpbEwZc98ZedHZG
  hcptf vcsevent read -id=ve-DJpbEwZc98ZedHZG -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the vcsevent read command
func (c *VCSEventReadCommand) Synopsis() string {
	return "Show details of a specific VCS event"
}
