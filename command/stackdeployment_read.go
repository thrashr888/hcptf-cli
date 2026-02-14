package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// StackDeploymentReadCommand is a command to read stack deployment details
type StackDeploymentReadCommand struct {
	Meta
	deploymentRunID string
	format          string
}

// Run executes the stack deployment read command
func (c *StackDeploymentReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackdeployment read")
	flags.StringVar(&c.deploymentRunID, "id", "", "Stack deployment run ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.deploymentRunID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read stack deployment run
	deploymentRun, err := client.StackDeploymentRuns.Read(client.Context(), c.deploymentRunID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading stack deployment: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":        deploymentRun.ID,
		"Status":    string(deploymentRun.Status),
		"CreatedAt": deploymentRun.CreatedAt,
		"UpdatedAt": deploymentRun.UpdatedAt,
	}

	// Add related resource information
	if deploymentRun.StackDeploymentGroup != nil {
		data["DeploymentGroupID"] = deploymentRun.StackDeploymentGroup.ID
		data["DeploymentGroupName"] = deploymentRun.StackDeploymentGroup.Name
	}

	// List deployment steps
	steps, err := client.StackDeploymentSteps.List(client.Context(), c.deploymentRunID, &tfe.StackDeploymentStepsListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 10,
		},
	})
	if err == nil && len(steps.Items) > 0 {
		data["TotalSteps"] = len(steps.Items)

		// Count steps by status
		statusCounts := make(map[string]int)
		for _, step := range steps.Items {
			statusCounts[string(step.Status)]++
		}

		stepStatus := []string{}
		for status, count := range statusCounts {
			stepStatus = append(stepStatus, fmt.Sprintf("%s: %d", status, count))
		}
		data["StepsStatus"] = strings.Join(stepStatus, ", ")
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack deployment read command
func (c *StackDeploymentReadCommand) Help() string {
	helpText := `
Usage: hcptf stackdeployment read [options]

  Read stack deployment details including status and progress.

Options:

  -id=<run-id>      Stack deployment run ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackdeployment read -id=sdr-abc123
  hcptf stackdeployment read -id=sdr-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack deployment read command
func (c *StackDeploymentReadCommand) Synopsis() string {
	return "Read stack deployment details and status"
}
