package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// StackCreateCommand is a command to create a stack
type StackCreateCommand struct {
	Meta
	name               string
	description        string
	projectID          string
	vcsIdentifier      string
	vcsBranch          string
	oauthTokenID       string
	serviceProvider    string
	speculativeEnabled bool
	format             string
}

// Run executes the stack create command
func (c *StackCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stack create")
	flags.StringVar(&c.name, "name", "", "Stack name (required)")
	flags.StringVar(&c.description, "description", "", "Stack description")
	flags.StringVar(&c.projectID, "project-id", "", "Project ID (required)")
	flags.StringVar(&c.vcsIdentifier, "vcs-identifier", "", "VCS repository identifier (org/repo)")
	flags.StringVar(&c.vcsBranch, "vcs-branch", "", "VCS branch (defaults to repo default branch)")
	flags.StringVar(&c.oauthTokenID, "oauth-token-id", "", "OAuth token ID for VCS connection")
	flags.StringVar(&c.serviceProvider, "service-provider", "github", "VCS service provider")
	flags.BoolVar(&c.speculativeEnabled, "speculative-enabled", false, "Enable speculative plans on PRs")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.projectID == "" {
		c.Ui.Error("Error: -project-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create stack options
	options := tfe.StackCreateOptions{
		Name:               c.name,
		SpeculativeEnabled: tfe.Bool(c.speculativeEnabled),
		Project: &tfe.Project{
			ID: c.projectID,
		},
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	// Configure VCS repo if provided
	if c.vcsIdentifier != "" {
		vcsRepo := &tfe.StackVCSRepoOptions{
			Identifier: c.vcsIdentifier,
		}

		if c.vcsBranch != "" {
			vcsRepo.Branch = c.vcsBranch
		}

		if c.oauthTokenID != "" {
			vcsRepo.OAuthTokenID = c.oauthTokenID
		}

		options.VCSRepo = vcsRepo
	}

	// Create stack
	stack, err := client.Stacks.Create(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating stack: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Stack '%s' created successfully", stack.Name))

	// Show stack details
	data := map[string]interface{}{
		"ID":                 stack.ID,
		"Name":               stack.Name,
		"Description":        stack.Description,
		"SpeculativeEnabled": stack.SpeculativeEnabled,
		"CreatedAt":          stack.CreatedAt,
	}

	if stack.Project != nil {
		data["Project"] = stack.Project.ID
	}

	if stack.VCSRepo != nil {
		data["VCSIdentifier"] = stack.VCSRepo.Identifier
		if stack.VCSRepo.Branch != "" {
			data["VCSBranch"] = stack.VCSRepo.Branch
		}
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack create command
func (c *StackCreateCommand) Help() string {
	helpText := `
Usage: hcptf stack create [options]

  Create a new stack. Stacks enable orchestrating deployments across multiple
  configurations and workspaces.

Options:

  -name=<name>                Stack name (required)
  -project-id=<id>            Project ID (required)
  -description=<text>         Stack description
  -vcs-identifier=<org/repo>  VCS repository identifier
  -vcs-branch=<branch>        VCS branch (defaults to repo default)
  -oauth-token-id=<id>        OAuth token ID for VCS connection
  -service-provider=<type>    VCS service provider (default: github)
  -speculative-enabled        Enable speculative plans on PRs
  -output=<format>            Output format: table (default) or json

Example:

  hcptf stack create -name=my-stack -project-id=prj-abc123
  hcptf stack create -name=infra-stack -project-id=prj-abc123 \
    -vcs-identifier=myorg/myrepo -oauth-token-id=ot-xyz789 \
    -description="Infrastructure stack"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack create command
func (c *StackCreateCommand) Synopsis() string {
	return "Create a new stack"
}
