package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// WorkspaceTagAddCommand is a command to add tags to a workspace
type WorkspaceTagAddCommand struct {
	Meta
	workspaceID string
	tags        string
}

// Run executes the workspacetag add command
func (c *WorkspaceTagAddCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspacetag add")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.workspaceID, "id", "", "Workspace ID (alias)")
	flags.StringVar(&c.tags, "tags", "", "Comma-separated list of tag names to add (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.workspaceID == "" {
		c.Ui.Error("Error: -workspace-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.tags == "" {
		c.Ui.Error("Error: -tags flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Parse tag names
	tagNames := strings.Split(c.tags, ",")
	var tagItems []*tfe.Tag
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name != "" {
			tagItems = append(tagItems, &tfe.Tag{
				Name: name,
			})
		}
	}

	if len(tagItems) == 0 {
		c.Ui.Error("Error: no valid tag names provided")
		return 1
	}

	// Add tags to workspace
	options := tfe.WorkspaceAddTagsOptions{
		Tags: tagItems,
	}

	err = client.Workspaces.AddTags(client.Context(), c.workspaceID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error adding tags to workspace: %s", err))
		return 1
	}

	if len(tagItems) == 1 {
		c.Ui.Output(fmt.Sprintf("Successfully added tag '%s' to workspace", tagItems[0].Name))
	} else {
		tagNameList := make([]string, len(tagItems))
		for i, tag := range tagItems {
			tagNameList[i] = tag.Name
		}
		c.Ui.Output(fmt.Sprintf("Successfully added %d tags to workspace: %s", len(tagItems), strings.Join(tagNameList, ", ")))
	}
	return 0
}

// Help returns help text for the workspacetag add command
func (c *WorkspaceTagAddCommand) Help() string {
	helpText := `
Usage: hcptf workspacetag add [options]

  Add organization tags to a workspace. Tags help categorize and organize
  workspaces. The tags must already exist in the organization - use the
  'organizationtag' commands to manage organization tags.

  You can add multiple tags by providing a comma-separated list of tag names.

Options:

  -workspace-id=<id>   Workspace ID (required)
  -id=<id>            Alias for -workspace-id
  -tags=<names>       Comma-separated list of tag names to add (required)

Example:

  hcptf workspacetag add -workspace-id=ws-ABC123 -tags=production
  hcptf workspacetag add -id=ws-ABC123 -tags=production,us-west-2,team-a
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspacetag add command
func (c *WorkspaceTagAddCommand) Synopsis() string {
	return "Add tags to a workspace"
}
