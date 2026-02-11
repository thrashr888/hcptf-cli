package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VCSEventListCommand lists VCS events for an organization
type VCSEventListCommand struct {
	Meta
	organization string
	format       string
	from         string
	to           string
	oauthClient  string
	level        string
}

// VCSEvent represents a VCS event
type VCSEvent struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes VCSEventAttributes   `json:"attributes"`
	Relationships VCSEventRelationships `json:"relationships,omitempty"`
}

// VCSEventAttributes contains VCS event details
type VCSEventAttributes struct {
	CreatedAt       string  `json:"created-at"`
	Level           string  `json:"level"`
	Message         string  `json:"message"`
	OrganizationID  string  `json:"organization-id"`
	SuggestedAction *string `json:"suggested_action,omitempty"`
}

// VCSEventRelationships contains related resources
type VCSEventRelationships struct {
	OAuthClient *struct {
		Data *struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	} `json:"oauth-client,omitempty"`
	OAuthToken *struct {
		Data *struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	} `json:"oauth-token,omitempty"`
}

// VCSEventListResponse represents the API response
type VCSEventListResponse struct {
	Data []VCSEvent `json:"data"`
	Meta struct {
		Pagination struct {
			CurrentPage int `json:"current-page"`
			PageSize    int `json:"page-size"`
			TotalPages  int `json:"total-pages"`
			TotalCount  int `json:"total-count"`
		} `json:"pagination"`
	} `json:"meta"`
}

// Run executes the vcsevent list command
func (c *VCSEventListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vcsevent list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")
	flags.StringVar(&c.from, "from", "", "Start time (RFC3339 format in UTC, e.g., 2021-02-02T14:09:00Z)")
	flags.StringVar(&c.to, "to", "", "End time (RFC3339 format in UTC, e.g., 2021-02-12T14:09:59Z)")
	flags.StringVar(&c.oauthClient, "oauth-client", "", "Filter by OAuth client external ID")
	flags.StringVar(&c.level, "level", "", "Filter by level: info or error")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate level if provided
	if c.level != "" && c.level != "info" && c.level != "error" {
		c.Ui.Error("Error: -level must be either 'info' or 'error'")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build API URL with query parameters
	apiURL := fmt.Sprintf("%s/api/v2/organizations/%s/vcs-events", client.GetAddress(), c.organization)

	// Build query parameters
	var queryParams []string
	if c.from != "" {
		queryParams = append(queryParams, fmt.Sprintf("filter[from]=%s", c.from))
	}
	if c.to != "" {
		queryParams = append(queryParams, fmt.Sprintf("filter[to]=%s", c.to))
	}
	if c.oauthClient != "" {
		queryParams = append(queryParams, fmt.Sprintf("filter[oauth_client_external_ids]=%s", c.oauthClient))
	}
	if c.level != "" {
		queryParams = append(queryParams, fmt.Sprintf("filter[levels]=%s", c.level))
	}

	if len(queryParams) > 0 {
		apiURL = apiURL + "?" + strings.Join(queryParams, "&")
	}

	req, err := http.NewRequestWithContext(client.Context(), "GET", apiURL, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating request: %s", err))
		return 1
	}

	// Get token from config for authorization
	u := client.BaseURL()
	cfg, err := c.Meta.Config()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error loading config: %s", err))
		return 1
	}
	token := cfg.GetToken(u.Hostname())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	httpClient := &http.Client{}
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
			c.Ui.Error("\nNote: VCS Events may not be available for your organization.")
			c.Ui.Error("The VCS Events API is in beta and currently only supports GitLab.com")
			c.Ui.Error("connections established after December 2020.")
		}
		return 1
	}

	var vcsEvents VCSEventListResponse
	if err := json.Unmarshal(body, &vcsEvents); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(vcsEvents.Data) == 0 {
		c.Ui.Output("No VCS events found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Level", "Message", "Created At"}
	var rows [][]string

	for _, event := range vcsEvents.Data {
		message := event.Attributes.Message
		if len(message) > 60 {
			message = message[:57] + "..."
		}

		rows = append(rows, []string{
			event.ID,
			event.Attributes.Level,
			message,
			event.Attributes.CreatedAt,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the vcsevent list command
func (c *VCSEventListCommand) Help() string {
	helpText := `
Usage: hcptf vcsevent list [options]

  List VCS (version control system) events for an organization.

  VCS events describe changes and actions related to VCS integration,
  helping debug webhook deliveries, OAuth token issues, and other
  VCS-related problems. Events are stored for 10 days.

  Note: The VCS Events API is in beta and currently only supports
  GitLab.com connections established after December 2020.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json
  -from=<timestamp>    Start time (RFC3339 format in UTC, e.g., 2021-02-02T14:09:00Z)
  -to=<timestamp>      End time (RFC3339 format in UTC, defaults to now)
  -oauth-client=<id>   Filter by OAuth client external ID
  -level=<level>       Filter by level: info or error

Example:

  hcptf vcsevent list -org=my-org
  hcptf vcsevent list -org=my-org -level=error
  hcptf vcsevent list -org=my-org -from=2021-02-02T14:09:00Z -to=2021-02-12T14:09:59Z
  hcptf vcsevent list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the vcsevent list command
func (c *VCSEventListCommand) Synopsis() string {
	return "List VCS events for an organization"
}
