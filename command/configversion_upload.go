package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ConfigVersionUploadCommand is a command to upload configuration files
type ConfigVersionUploadCommand struct {
	Meta
	configVersionID string
	path            string
	configVerSvc    configVersionReader
	uploadSvc       configVersionUploader
}

// Run executes the configversion upload command
func (c *ConfigVersionUploadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("configversion upload")
	flags.StringVar(&c.configVersionID, "id", "", "Configuration version ID (required)")
	flags.StringVar(&c.path, "path", "", "Path to configuration directory or tar.gz file (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.configVersionID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.path == "" {
		c.Ui.Error("Error: -path flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Verify path exists
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		c.Ui.Error(fmt.Sprintf("Error: path does not exist: %s", c.path))
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get configuration version to retrieve upload URL
	configVersion, err := c.configVersionService(client).Read(client.Context(), c.configVersionID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading configuration version: %s", err))
		return 1
	}

	if configVersion.UploadURL == "" {
		c.Ui.Error("Error: configuration version does not have an upload URL")
		c.Ui.Error("The upload URL is only available immediately after creation.")
		return 1
	}

	// Upload the configuration
	c.Ui.Output(fmt.Sprintf("Uploading configuration from: %s", c.path))
	err = c.uploadService(client).Upload(client.Context(), configVersion.UploadURL, c.path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error uploading configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully uploaded configuration to version: %s", c.configVersionID))
	return 0
}

func (c *ConfigVersionUploadCommand) configVersionService(client *client.Client) configVersionReader {
	if c.configVerSvc != nil {
		return c.configVerSvc
	}
	return client.ConfigurationVersions
}

func (c *ConfigVersionUploadCommand) uploadService(client *client.Client) configVersionUploader {
	if c.uploadSvc != nil {
		return c.uploadSvc
	}
	return client.ConfigurationVersions
}

// Help returns help text for the configversion upload command
func (c *ConfigVersionUploadCommand) Help() string {
	helpText := `
Usage: hcptf configversion upload [options]

  Upload configuration files to a configuration version.

  The path can be either a directory containing Terraform configuration
  files or a tar.gz archive. If a directory is provided, it will be
  automatically archived before upload.

Options:

  -id=<config-id>   Configuration version ID (required)
  -path=<path>      Path to configuration directory or tar.gz file (required)

Example:

  hcptf configversion upload -id=cv-abc123 -path=./terraform
  hcptf configversion upload -id=cv-abc123 -path=./config.tar.gz
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the configversion upload command
func (c *ConfigVersionUploadCommand) Synopsis() string {
	return "Upload configuration files"
}
