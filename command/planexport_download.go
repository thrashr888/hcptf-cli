package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PlanExportDownloadCommand is a command to download exported plan data
type PlanExportDownloadCommand struct {
	Meta
	planExportID  string
	outputPath    string
	planExportSvc planExportDownloader
}

// Run executes the planexport download command
func (c *PlanExportDownloadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("planexport download")
	flags.StringVar(&c.planExportID, "id", "", "Plan export ID (required)")
	flags.StringVar(&c.outputPath, "path", "", "Output file path (default: <export-id>.tar.gz)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.planExportID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Set default output path if not specified
	if c.outputPath == "" {
		c.outputPath = fmt.Sprintf("%s.tar.gz", c.planExportID)
	}

	// Check if output file already exists
	if _, err := os.Stat(c.outputPath); err == nil {
		c.Ui.Error(fmt.Sprintf("Error: output file already exists: %s", c.outputPath))
		c.Ui.Error("Please specify a different path or remove the existing file")
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Download plan export data
	c.Ui.Output(fmt.Sprintf("Downloading plan export: %s", c.planExportID))
	data, err := c.planExportService(client).Download(client.Context(), c.planExportID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error downloading plan export: %s", err))
		return 1
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(c.outputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			c.Ui.Error(fmt.Sprintf("Error creating output directory: %s", err))
			return 1
		}
	}

	// Write data to file
	if err := os.WriteFile(c.outputPath, data, 0644); err != nil {
		c.Ui.Error(fmt.Sprintf("Error writing file: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully downloaded plan export to: %s", c.outputPath))
	c.Ui.Output(fmt.Sprintf("File size: %d bytes", len(data)))
	return 0
}

func (c *PlanExportDownloadCommand) planExportService(client *client.Client) planExportDownloader {
	if c.planExportSvc != nil {
		return c.planExportSvc
	}
	return client.PlanExports
}

// Help returns help text for the planexport download command
func (c *PlanExportDownloadCommand) Help() string {
	helpText := `
Usage: hcptf planexport download [options]

  Download exported plan data to a local file.

  The plan export must be in 'finished' status before it can be downloaded.
  Use 'planexport read' to check the export status first.

  The exported data is saved as a .tar.gz archive containing the plan
  data in the requested format (e.g., Sentinel mock data).

Options:

  -id=<export-id>     Plan export ID (required)
  -path=<file>        Output file path (default: <export-id>.tar.gz)

Example:

  hcptf planexport download -id=pe-abc123
  hcptf planexport download -id=pe-abc123 -path=./exports/plan-export.tar.gz
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the planexport download command
func (c *PlanExportDownloadCommand) Synopsis() string {
	return "Download exported plan data"
}
