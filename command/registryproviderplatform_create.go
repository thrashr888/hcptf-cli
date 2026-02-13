package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderPlatformCreateCommand struct {
	Meta
	organization string
	name         string
	namespace    string
	version      string
	os           string
	arch         string
	shasum       string
	filename     string
	format       string
}

// Run executes the registry provider platform create command
func (c *RegistryProviderPlatformCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryproviderplatform create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.version, "version", "", "Version string (required)")
	flags.StringVar(&c.os, "os", "", "Operating system, e.g., linux, darwin, windows (required)")
	flags.StringVar(&c.arch, "arch", "", "Architecture, e.g., amd64, arm64, 386 (required)")
	flags.StringVar(&c.shasum, "shasum", "", "SHA256 checksum (required)")
	flags.StringVar(&c.filename, "filename", "", "Binary filename (required)")
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

	if c.version == "" {
		c.Ui.Error("Error: -version flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.os == "" {
		c.Ui.Error("Error: -os flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.arch == "" {
		c.Ui.Error("Error: -arch flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.shasum == "" {
		c.Ui.Error("Error: -shasum flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.filename == "" {
		c.Ui.Error("Error: -filename flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Default namespace to organization
	if c.namespace == "" {
		c.namespace = c.organization
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create provider platform
	versionID := tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			RegistryName: tfe.PrivateRegistry,
			Namespace:    c.namespace,
			Name:         c.name,
		},
		Version: c.version,
	}

	options := tfe.RegistryProviderPlatformCreateOptions{
		OS:       c.os,
		Arch:     c.arch,
		Shasum:   c.shasum,
		Filename: c.filename,
	}

	platform, err := client.RegistryProviderPlatforms.Create(client.Context(), versionID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating provider platform: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Provider platform '%s/%s' created successfully", c.os, c.arch))

	if uploadURL, ok := platform.Links["provider-binary-upload"]; ok {
		c.Ui.Output(fmt.Sprintf("\nUpload the provider binary to: %s", uploadURL))
		c.Ui.Output("\nExample:")
		c.Ui.Output(fmt.Sprintf("  curl --header \"Authorization: Bearer $TOKEN\" --header \"Content-Type: application/octet-stream\" --request PUT --data-binary @%s %s", c.filename, uploadURL))
	}

	// Show platform details
	data := map[string]interface{}{
		"ID":       platform.ID,
		"OS":       platform.OS,
		"Arch":     platform.Arch,
		"Filename": platform.Filename,
		"Shasum":   platform.Shasum,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the registry provider platform create command
func (c *RegistryProviderPlatformCreateCommand) Help() string {
	helpText := `
Usage: hcptf registryproviderplatform create [options]

  Create a new platform binary for a private registry provider version.

  After creating the platform, you must upload the provider binary
  to the provided upload URL.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -version=<semver>    Version string (required)
  -os=<os>             Operating system: linux, darwin, windows, etc. (required)
  -arch=<arch>         Architecture: amd64, arm64, 386, etc. (required)
  -shasum=<sha256>     SHA256 checksum of the binary file (required)
  -filename=<name>     Binary filename (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryproviderplatform create -org=my-org -name=aws -version=3.1.1 \
    -os=linux -arch=amd64 \
    -shasum=8f69533bc8afc227b40d15116358f91505bb638ce5919712fbb38a2dec1bba38 \
    -filename=terraform-provider-aws_3.1.1_linux_amd64.zip
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider platform create command
func (c *RegistryProviderPlatformCreateCommand) Synopsis() string {
	return "Create a platform binary for a private registry provider version"
}
