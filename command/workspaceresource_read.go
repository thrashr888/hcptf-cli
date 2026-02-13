package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// WorkspaceResourceReadCommand is a command to read workspace resource details
type WorkspaceResourceReadCommand struct {
	Meta
	workspaceID string
	resourceID  string
	format      string
}

// Run executes the workspace resource read command
func (c *WorkspaceResourceReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspaceresource read")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.resourceID, "id", "", "Resource ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.workspaceID == "" {
		c.Ui.Error("Error: -workspace-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.resourceID == "" {
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

	// List workspace resources and find the specific one
	// Note: The API only provides a List endpoint, not a Read endpoint
	resources, err := client.WorkspaceResources.List(client.Context(), c.workspaceID, &tfe.WorkspaceResourceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing workspace resources: %s", err))
		return 1
	}

	// Find the specific resource
	var resource *tfe.WorkspaceResource
	for _, r := range resources.Items {
		if r.ID == c.resourceID {
			resource = r
			break
		}
	}

	if resource == nil {
		c.Ui.Error(fmt.Sprintf("Error: Resource with ID %s not found in workspace", c.resourceID))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                       resource.ID,
		"Address":                  resource.Address,
		"Name":                     resource.Name,
		"Type":                     resource.ProviderType,
		"Provider":                 resource.Provider,
		"Module":                   resource.Module,
		"CreatedAt":                resource.CreatedAt,
		"UpdatedAt":                resource.UpdatedAt,
		"ModifiedByStateVersionID": resource.ModifiedByStateVersionID,
	}

	if resource.NameIndex != nil {
		data["NameIndex"] = *resource.NameIndex
	} else {
		data["NameIndex"] = "N/A"
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the workspace resource read command
func (c *WorkspaceResourceReadCommand) Help() string {
	helpText := `
Usage: hcptf workspaceresource read [options]

  Show details for a specific resource in workspace state.

Options:

  -workspace-id=<id>   Workspace ID (required, format: ws-xxx)
  -id=<resource-id>    Resource ID (required, format: wsr-xxx)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspaceresource read -workspace-id=ws-abc123 -id=wsr-xyz789
  hcptf workspaceresource read -workspace-id=ws-abc123 -id=wsr-xyz789 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace resource read command
func (c *WorkspaceResourceReadCommand) Synopsis() string {
	return "Show resource details"
}
