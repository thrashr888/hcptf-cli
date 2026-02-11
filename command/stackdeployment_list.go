package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// StackDeploymentListCommand is a command to list stack deployments
type StackDeploymentListCommand struct {
	Meta
	stackID string
	format  string
}

// Run executes the stack deployment list command
func (c *StackDeploymentListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackdeployment list")
	flags.StringVar(&c.stackID, "stack-id", "", "Stack ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.stackID == "" {
		c.Ui.Error("Error: -stack-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List stack deployments
	deployments, err := client.StackDeployments.List(client.Context(), c.stackID, &tfe.StackDeploymentListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing stack deployments: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(deployments.Items) == 0 {
		c.Ui.Output("No stack deployments found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Latest Run Status"}
	var rows [][]string

	for _, deployment := range deployments.Items {
		runStatus := "N/A"
		if deployment.LatestDeploymentRun != nil {
			runStatus = string(deployment.LatestDeploymentRun.Status)
		}

		rows = append(rows, []string{
			deployment.ID,
			deployment.Name,
			runStatus,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the stack deployment list command
func (c *StackDeploymentListCommand) Help() string {
	helpText := `
Usage: hcptf stackdeployment list [options]

  List deployments for a stack. Deployments represent the execution of
  stack configurations across your infrastructure.

Options:

  -stack-id=<id>    Stack ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackdeployment list -stack-id=st-abc123
  hcptf stackdeployment list -stack-id=st-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack deployment list command
func (c *StackDeploymentListCommand) Synopsis() string {
	return "List deployments for a stack"
}
