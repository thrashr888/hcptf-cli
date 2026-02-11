package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type RegistryModuleCreateVersionCommand struct {
	Meta
	organization      string
	name              string
	provider          string
	namespace         string
	registryName      string
	version           string
	commitSHA         string
	format            string
	registryModuleSvc registryModuleVersionCreator
}

// Run executes the registry module create-version command
func (c *RegistryModuleCreateVersionCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registrymodule create-version")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Module name (required)")
	flags.StringVar(&c.provider, "provider", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.registryName, "registry-name", "private", "Registry name: must be 'private'")
	flags.StringVar(&c.version, "version", "", "Version string, e.g., 1.0.0 (required)")
	flags.StringVar(&c.commitSHA, "commit-sha", "", "Commit SHA for the version")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.provider == "" {
		c.Ui.Error("Error: -provider flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.version == "" {
		c.Ui.Error("Error: -version flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Default namespace to organization for private modules
	if c.namespace == "" {
		c.namespace = c.organization
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create module version
	moduleID := tfe.RegistryModuleID{
		Organization: c.organization,
		Name:         c.name,
		Provider:     c.provider,
		Namespace:    c.namespace,
		RegistryName: tfe.RegistryName(c.registryName),
	}

	options := tfe.RegistryModuleCreateVersionOptions{
		Version: tfe.String(c.version),
	}

	if c.commitSHA != "" {
		options.CommitSHA = tfe.String(c.commitSHA)
	}

	moduleVersion, err := c.registryModuleService(client).CreateVersion(client.Context(), moduleID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating module version: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Module version '%s' created successfully", moduleVersion.Version))
	c.Ui.Output(fmt.Sprintf("\nUpload URL: %s", moduleVersion.Links["upload"]))
	c.Ui.Output("\nUse the upload URL to upload your module tarball:")
	c.Ui.Output(fmt.Sprintf("  curl --header \"Authorization: Bearer $TOKEN\" --header \"Content-Type: application/octet-stream\" --request PUT --data-binary @module.tar.gz %s", moduleVersion.Links["upload"]))

	// Show version details
	data := map[string]interface{}{
		"ID":        moduleVersion.ID,
		"Version":   moduleVersion.Version,
		"Status":    moduleVersion.Status,
		"Source":    moduleVersion.Source,
		"CreatedAt": moduleVersion.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RegistryModuleCreateVersionCommand) registryModuleService(client *client.Client) registryModuleVersionCreator {
	if c.registryModuleSvc != nil {
		return c.registryModuleSvc
	}
	return client.RegistryModules
}

// Help returns help text for the registry module create-version command
func (c *RegistryModuleCreateVersionCommand) Help() string {
	helpText := `
Usage: hcptf registrymodule create-version [options]

  Create a new version for a private registry module.

  After creating the version, you must upload the module tarball
  to the provided upload URL.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Module name (required)
  -provider=<name>     Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -registry-name=<val> Registry name: must be 'private'
  -version=<semver>    Version string, e.g., 1.0.0 (required)
  -commit-sha=<sha>    Commit SHA for the version
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registrymodule create-version -org=my-org -name=vpc -provider=aws -version=1.2.3
  hcptf registrymodule create-version -org=my-org -name=vpc -provider=aws -version=1.2.3 -commit-sha=abc123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry module create-version command
func (c *RegistryModuleCreateVersionCommand) Synopsis() string {
	return "Create a new version for a private registry module"
}
