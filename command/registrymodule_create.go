package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type RegistryModuleCreateCommand struct {
	Meta
	organization      string
	name              string
	provider          string
	registryName      string
	noCode            bool
	format            string
	registryModuleSvc registryModuleCreator
}

// Run executes the registry module create command
func (c *RegistryModuleCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registrymodule create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Module name (required)")
	flags.StringVar(&c.provider, "provider", "", "Provider name (required)")
	flags.StringVar(&c.registryName, "registry-name", "private", "Registry name: public or private (default: private)")
	flags.BoolVar(&c.noCode, "no-code", false, "Enable no-code publishing workflow")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create registry module
	options := tfe.RegistryModuleCreateOptions{
		Name:         tfe.String(c.name),
		Provider:     tfe.String(c.provider),
		RegistryName: tfe.RegistryName(c.registryName),
	}

	if c.noCode {
		options.NoCode = tfe.Bool(c.noCode)
	}

	module, err := c.registryModuleService(client).Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating registry module: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Registry module '%s/%s' created successfully", module.Name, module.Provider))

	// Show module details
	data := map[string]interface{}{
		"ID":           module.ID,
		"Name":         module.Name,
		"Namespace":    module.Namespace,
		"Provider":     module.Provider,
		"RegistryName": string(module.RegistryName),
		"Status":       module.Status,
		"CreatedAt":    module.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RegistryModuleCreateCommand) registryModuleService(client *client.Client) registryModuleCreator {
	if c.registryModuleSvc != nil {
		return c.registryModuleSvc
	}
	return client.RegistryModules
}

// Help returns help text for the registry module create command
func (c *RegistryModuleCreateCommand) Help() string {
	helpText := `
Usage: hcptf registrymodule create [options]

  Create a new private registry module without a VCS connection.

  After creating a module, you must create and upload versions using
  the 'registrymodule create-version' command.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Module name (required)
  -provider=<name>     Provider name, e.g., aws, azure, gcp (required)
  -registry-name=<val> Registry name: public or private (default: private)
  -no-code             Enable no-code publishing workflow
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registrymodule create -org=my-org -name=vpc -provider=aws
  hcptf registrymodule create -org=my-org -name=network -provider=azure -no-code
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry module create command
func (c *RegistryModuleCreateCommand) Synopsis() string {
	return "Create a new private registry module"
}
