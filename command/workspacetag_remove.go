package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// WorkspaceTagRemoveCommand is a command to remove tags from a workspace
type WorkspaceTagRemoveCommand struct {
	Meta
	workspaceID string
	tags        string
}

// Run executes the workspacetag remove command
func (c *WorkspaceTagRemoveCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspacetag remove")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.workspaceID, "id", "", "Workspace ID (alias)")
	flags.StringVar(&c.tags, "tags", "", "Comma-separated list of tag names to remove (required)")

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

	// Remove tags from workspace
	options := tfe.WorkspaceRemoveTagsOptions{
		Tags: tagItems,
	}

	err = client.Workspaces.RemoveTags(client.Context(), c.workspaceID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing tags from workspace: %s", err))
		return 1
	}

	if len(tagItems) == 1 {
		c.Ui.Output(fmt.Sprintf("Successfully removed tag '%s' from workspace", tagItems[0].Name))
	} else {
		tagNameList := make([]string, len(tagItems))
		for i, tag := range tagItems {
			tagNameList[i] = tag.Name
		}
		c.Ui.Output(fmt.Sprintf("Successfully removed %d tags from workspace: %s", len(tagItems), strings.Join(tagNameList, ", ")))
	}
	return 0
}

// Help returns help text for the workspacetag remove command
func (c *WorkspaceTagRemoveCommand) Help() string {
	helpText := `
Usage: hcptf workspacetag remove [options]

  Remove organization tags from a workspace. This only removes the tags
  from the workspace - the tags will still exist in the organization and
  can be applied to other workspaces.

  You can remove multiple tags by providing a comma-separated list of tag names.

Options:

  -workspace-id=<id>   Workspace ID (required)
  -id=<id>            Alias for -workspace-id
  -tags=<names>       Comma-separated list of tag names to remove (required)

Example:

  hcptf workspacetag remove -workspace-id=ws-ABC123 -tags=production
  hcptf workspacetag remove -id=ws-ABC123 -tags=production,us-west-2,team-a
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspacetag remove command
func (c *WorkspaceTagRemoveCommand) Synopsis() string {
	return "Remove tags from a workspace"
}
