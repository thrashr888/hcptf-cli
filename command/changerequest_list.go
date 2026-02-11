package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ChangeRequestListCommand lists change requests for a workspace
type ChangeRequestListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

// ChangeRequest represents a workspace change request
type ChangeRequest struct {
	ID            string                     `json:"id"`
	Type          string                     `json:"type"`
	Attributes    ChangeRequestAttributes    `json:"attributes"`
	Relationships ChangeRequestRelationships `json:"relationships"`
}

// ChangeRequestAttributes contains change request details
type ChangeRequestAttributes struct {
	Subject    string  `json:"subject"`
	Message    string  `json:"message"`
	ArchivedBy *string `json:"archived-by"`
	ArchivedAt *string `json:"archived-at"`
	CreatedAt  string  `json:"created-at"`
	UpdatedAt  string  `json:"updated-at"`
}

// ChangeRequestRelationships contains related workspace
type ChangeRequestRelationships struct {
	Workspace struct {
		Data struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	} `json:"workspace"`
}

// ChangeRequestListResponse represents the API response
type ChangeRequestListResponse struct {
	Data []ChangeRequest `json:"data"`
	Meta struct {
		Pagination struct {
			CurrentPage int `json:"current-page"`
			PageSize    int `json:"page-size"`
			TotalPages  int `json:"total-pages"`
			TotalCount  int `json:"total-count"`
		} `json:"pagination"`
	} `json:"meta"`
}

// Run executes the changerequest list command
func (c *ChangeRequestListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("changerequest list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
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

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace to obtain its ID
	workspace, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Make direct API call to list change requests
	apiURL := fmt.Sprintf("%s/api/v2/workspaces/%s/change-requests", client.GetAddress(), workspace.ID)

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
			c.Ui.Error("\nNote: Change Requests may not be available in your HCP Terraform plan.")
			c.Ui.Error("This feature requires HCP Terraform Plus or Enterprise.")
		}
		return 1
	}

	var changeRequests ChangeRequestListResponse
	if err := json.Unmarshal(body, &changeRequests); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(changeRequests.Data) == 0 {
		c.Ui.Output("No change requests found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Subject", "Status", "Created At"}
	var rows [][]string

	for _, cr := range changeRequests.Data {
		status := "Open"
		if cr.Attributes.ArchivedAt != nil {
			status = "Archived"
		}

		subject := cr.Attributes.Subject
		if len(subject) > 60 {
			subject = subject[:57] + "..."
		}

		rows = append(rows, []string{
			cr.ID,
			subject,
			status,
			cr.Attributes.CreatedAt,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the changerequest list command
func (c *ChangeRequestListCommand) Help() string {
	helpText := `
Usage: hcptf changerequest list [options]

  List change requests for a workspace.

  Change requests are used to keep track of workspace to-dos, helping
  teams manage compliance and best practices. They can trigger team
  notifications when created.

  Note: This feature requires HCP Terraform Plus or Enterprise.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf changerequest list -org=my-org -workspace=my-workspace
  hcptf changerequest list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the changerequest list command
func (c *ChangeRequestListCommand) Synopsis() string {
	return "List change requests for a workspace"
}
