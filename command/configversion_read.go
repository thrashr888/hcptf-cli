package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ConfigVersionReadCommand is a command to read configuration version details
type ConfigVersionReadCommand struct {
	Meta
	configVersionID string
	runID           string
	format          string
	configVerSvc    configVersionReader
	runSvc          runReader
}

// Run executes the configversion read command
func (c *ConfigVersionReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("configversion read")
	flags.StringVar(&c.configVersionID, "id", "", "Configuration version ID or Run ID")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (alternative to -id)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags - need either configVersionID or runID
	id := c.configVersionID
	if id == "" {
		id = c.runID
	}
	if id == "" {
		c.Ui.Error("Error: -id or -run-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// If ID starts with "run-", get the config version ID from the run
	configVersionID := id
	if strings.HasPrefix(id, "run-") {
		run, err := c.runService(client).Read(client.Context(), id)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading run: %s", err))
			return 1
		}
		if run.ConfigurationVersion == nil {
			c.Ui.Error("Error: run has no configuration version")
			return 1
		}
		configVersionID = run.ConfigurationVersion.ID
	}

	// Read configuration version
	configVersion, err := c.configVersionService(client).Read(client.Context(), configVersionID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading configuration version: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":           configVersion.ID,
		"Status":       configVersion.Status,
		"Source":       configVersion.Source,
		"Speculative":  configVersion.Speculative,
		"Provisional":  configVersion.Provisional,
		"Error":        configVersion.Error,
		"ErrorMessage": configVersion.ErrorMessage,
	}

	if configVersion.UploadURL != "" {
		data["UploadURL"] = configVersion.UploadURL
	}

	// Fetch VCS ingress attributes if this is from a VCS source
	if configVersion.Source == "github" || configVersion.Source == "gitlab" || configVersion.Source == "bitbucket" || configVersion.Source == "ado" {
		if ingress := c.fetchIngressAttributes(client, configVersionID); ingress != nil {
			if ingress.Branch != "" {
				data["Branch"] = ingress.Branch
			}
			if ingress.CommitSHA != "" {
				data["CommitSHA"] = ingress.CommitSHA
			}
			if ingress.CommitURL != "" {
				data["CommitURL"] = ingress.CommitURL
			}
			if ingress.CompareURL != "" {
				data["CompareURL"] = ingress.CompareURL
			}
			if ingress.Identifier != "" {
				data["RepoIdentifier"] = ingress.Identifier
			}
		}
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ConfigVersionReadCommand) configVersionService(client *client.Client) configVersionReader {
	if c.configVerSvc != nil {
		return c.configVerSvc
	}
	return client.ConfigurationVersions
}

func (c *ConfigVersionReadCommand) runService(client *client.Client) runReader {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// IngressAttributes represents VCS ingress information
type IngressAttributes struct {
	Branch     string `json:"branch"`
	CommitSHA  string `json:"commit-sha"`
	CommitURL  string `json:"commit-url"`
	CompareURL string `json:"compare-url"`
	Identifier string `json:"identifier"`
}

// IngressAttributesResponse represents the API response for ingress attributes
type IngressAttributesResponse struct {
	Data struct {
		ID         string             `json:"id"`
		Type       string             `json:"type"`
		Attributes IngressAttributes  `json:"attributes"`
	} `json:"data"`
}

// fetchIngressAttributes fetches VCS ingress attributes for a configuration version
func (c *ConfigVersionReadCommand) fetchIngressAttributes(client *client.Client, configVersionID string) *IngressAttributes {
	apiURL := fmt.Sprintf("%s/api/v2/configuration-versions/%s/ingress-attributes", client.GetAddress(), configVersionID)

	req, err := http.NewRequestWithContext(client.Context(), "GET", apiURL, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+client.Token())
	req.Header.Set("Content-Type", "application/vnd.api+json")

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var response IngressAttributesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	return &response.Data.Attributes
}

// Help returns help text for the configversion read command
func (c *ConfigVersionReadCommand) Help() string {
	helpText := `
Usage: hcptf configversion read [options]

  Show configuration version details. You can provide either a configuration
  version ID or a run ID. If you provide a run ID, the command will automatically
  look up the associated configuration version.

Options:

  -id=<id>          Configuration version ID (cv-xxx) or Run ID (run-xxx) (required)
  -run-id=<id>      Run ID (alternative to -id)
  -output=<format>  Output format: table (default) or json

Examples:

  # Using configuration version ID
  hcptf configversion read -id=cv-abc123

  # Using run ID
  hcptf configversion read -id=run-xyz789
  hcptf configversion read -run-id=run-xyz789

  # URL-style
  hcptf my-org my-workspace runs run-xyz789 configversion
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the configversion read command
func (c *ConfigVersionReadCommand) Synopsis() string {
	return "Show configuration version details"
}
