package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type RegistryModuleListCommand struct {
	Meta
	organization      string
	format            string
	registryModuleSvc registryModuleLister
}

// Run executes the registry module list command
func (c *RegistryModuleListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registrymodule list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List registry modules
	modules, err := c.registryModuleService(client).List(client.Context(), c.organization, &tfe.RegistryModuleListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing registry modules: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(modules.Items) == 0 {
		c.Ui.Output("No registry modules found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Namespace", "Provider", "Registry", "Status"}
	var rows [][]string

	for _, mod := range modules.Items {
		rows = append(rows, []string{
			mod.ID,
			mod.Name,
			mod.Namespace,
			mod.Provider,
			string(mod.RegistryName),
			string(mod.Status),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *RegistryModuleListCommand) registryModuleService(client *client.Client) registryModuleLister {
	if c.registryModuleSvc != nil {
		return c.registryModuleSvc
	}
	return client.RegistryModules
}

// Help returns help text for the registry module list command
func (c *RegistryModuleListCommand) Help() string {
	helpText := `
Usage: hcptf registrymodule list [options]

  List private registry modules in an organization.

  The private registry enables organizations to publish and share
  internal Terraform modules within their organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registrymodule list -organization=my-org
  hcptf registrymodule list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry module list command
func (c *RegistryModuleListCommand) Synopsis() string {
	return "List private registry modules in an organization"
}
