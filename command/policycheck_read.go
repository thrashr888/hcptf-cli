package command

import (
	"fmt"
	"strings"
)

// PolicyCheckReadCommand is a command to read policy check details
type PolicyCheckReadCommand struct {
	Meta
	policyCheckID string
	format        string
}

// Run executes the policy check read command
func (c *PolicyCheckReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policycheck read")
	flags.StringVar(&c.policyCheckID, "id", "", "Policy Check ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policyCheckID == "" {
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

	// Read policy check
	policyCheck, err := client.PolicyChecks.Read(client.Context(), c.policyCheckID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy check: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":            policyCheck.ID,
		"Status":        string(policyCheck.Status),
		"Scope":         string(policyCheck.Scope),
		"IsOverridable": policyCheck.Actions.IsOverridable,
		"CanOverride":   policyCheck.Permissions.CanOverride,
	}

	if policyCheck.Result != nil {
		data["Passed"] = policyCheck.Result.Passed
		data["TotalFailed"] = policyCheck.Result.TotalFailed
		data["HardFailed"] = policyCheck.Result.HardFailed
		data["SoftFailed"] = policyCheck.Result.SoftFailed
		data["AdvisoryFailed"] = policyCheck.Result.AdvisoryFailed
		data["DurationMs"] = policyCheck.Result.Duration

		// Include detailed Sentinel policy results if available
		if policyCheck.Result.Sentinel != nil {
			data["Sentinel"] = policyCheck.Result.Sentinel
		}
	}

	if policyCheck.StatusTimestamps != nil {
		if !policyCheck.StatusTimestamps.QueuedAt.IsZero() {
			data["QueuedAt"] = policyCheck.StatusTimestamps.QueuedAt
		}
		if !policyCheck.StatusTimestamps.PassedAt.IsZero() {
			data["PassedAt"] = policyCheck.StatusTimestamps.PassedAt
		}
		if !policyCheck.StatusTimestamps.HardFailedAt.IsZero() {
			data["HardFailedAt"] = policyCheck.StatusTimestamps.HardFailedAt
		}
		if !policyCheck.StatusTimestamps.SoftFailedAt.IsZero() {
			data["SoftFailedAt"] = policyCheck.StatusTimestamps.SoftFailedAt
		}
		if !policyCheck.StatusTimestamps.ErroredAt.IsZero() {
			data["ErroredAt"] = policyCheck.StatusTimestamps.ErroredAt
		}
	}

	formatter.KeyValue(data)

	// Display detailed policy results if available and not in JSON mode
	if c.format != "json" && policyCheck.Result != nil && policyCheck.Result.Sentinel != nil {
		c.Ui.Output("\n=== POLICY RESULTS ===\n")

		// Parse Sentinel data structure
		sentinelData, ok := policyCheck.Result.Sentinel.(map[string]interface{})
		if !ok {
			return 0
		}

		data, ok := sentinelData["data"].(map[string]interface{})
		if !ok {
			return 0
		}

		// Iterate through policy sets (usually has empty key for the main set)
		for setName, setData := range data {
			setMap, ok := setData.(map[string]interface{})
			if !ok {
				continue
			}

			policies, ok := setMap["policies"].([]interface{})
			if !ok {
				continue
			}

			if setName != "" {
				c.Ui.Output(fmt.Sprintf("Policy Set: %s\n", setName))
			}

			// Display each policy result
			for _, policyData := range policies {
				policy, ok := policyData.(map[string]interface{})
				if !ok {
					continue
				}

				policyInfo, ok := policy["policy"].(map[string]interface{})
				if !ok {
					continue
				}

				name := policyInfo["name"]
				enforcementLevel := policyInfo["enforcement-level"]
				duration := policy["duration"]

				// Determine status symbol using safe type assertion
				resultBool, resultOk := policy["result"].(bool)

				statusSymbol := "✓"
				statusText := "PASSED"
				if resultOk && !resultBool {
					statusSymbol = "✗"
					statusText = "FAILED"
				}

				c.Ui.Output(fmt.Sprintf("%s %s", statusSymbol, statusText))
				c.Ui.Output(fmt.Sprintf("  Policy: %s", name))
				c.Ui.Output(fmt.Sprintf("  Enforcement: %s", enforcementLevel))
				c.Ui.Output(fmt.Sprintf("  Duration: %vms", duration))

				// Show error if present
				if err, ok := policy["error"]; ok && err != nil {
					c.Ui.Output(fmt.Sprintf("  Error: %v", err))
				}

				// Show trace information for failures
				if resultOk && !resultBool {
					trace, ok := policy["trace"].(map[string]interface{})
					if ok {
						// Show print output (policy messages)
						if printOutput, ok := trace["print"].(string); ok && printOutput != "" {
							c.Ui.Output(fmt.Sprintf("  Output:\n%s", indent(printOutput, "    ")))
						}

						// Show error message from trace
						if traceErr, ok := trace["error"]; ok && traceErr != nil && traceErr != "" {
							c.Ui.Output(fmt.Sprintf("  Message: %v", traceErr))
						}

						// Show rule results
						if rules, ok := trace["rules"].(map[string]interface{}); ok {
							for ruleName, ruleData := range rules {
								ruleMap, ok := ruleData.(map[string]interface{})
								if !ok {
									continue
								}
								if ruleValue, ok := ruleMap["value"].(bool); ok && !ruleValue {
									c.Ui.Output(fmt.Sprintf("  Rule '%s': %v", ruleName, ruleValue))
									if desc, ok := ruleMap["desc"].(string); ok && desc != "" {
										c.Ui.Output(fmt.Sprintf("    Description: %s", desc))
									}
								}
							}
						}
					}
				}

				c.Ui.Output("") // Blank line between policies
			}
		}
	}

	return 0
}

// indent adds a prefix to each line of text
func indent(text, prefix string) string {
	if text == "" {
		return text
	}
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// Help returns help text for the policy check read command
func (c *PolicyCheckReadCommand) Help() string {
	helpText := `
Usage: hcptf policy check read [options]

  Read policy check details and results.

Options:

  -id=<policy-check-id>  Policy Check ID (required)
  -output=<format>       Output format: table (default) or json

Example:

  hcptf policy check read -id=polchk-abc123
  hcptf policy check read -id=polchk-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy check read command
func (c *PolicyCheckReadCommand) Synopsis() string {
	return "Read policy check details"
}
