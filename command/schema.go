package command

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// SchemaCommand returns machine-readable schema for command flags.
type SchemaCommand struct {
	Meta
}

type schemaFlag struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Default     string   `json:"default,omitempty"`
}

type schemaOutput struct {
	Command string       `json:"command"`
	Flags   []schemaFlag `json:"flags"`
}

// Run executes the schema command.
func (c *SchemaCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Ui.Error("Error: schema requires a command path")
		c.Ui.Error(c.Help())
		return 1
	}

	path := strings.Join(args, " ")
	factory, ok := Commands(&c.Meta)[path]
	if !ok {
		c.Ui.Error(fmt.Sprintf("Error: unknown command %q", path))
		return 1
	}

	cmd, err := factory()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error loading command %q: %s", path, err))
		return 1
	}

	flags := parseFlagSchemaFromHelp(cmd.Help())
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].Name < flags[j].Name
	})

	formatter := c.Meta.NewFormatter("json")
	formatter.JSON(schemaOutput{Command: path, Flags: flags})
	return 0
}

func parseFlagSchemaFromHelp(helpText string) []schemaFlag {
	lines := strings.Split(helpText, "\n")
	flagNamePattern := regexp.MustCompile(`^\s*-([A-Za-z0-9_-]+)(?:(?:=)<[^>]+>)?\s+(.*)$`)
	aliasPattern := regexp.MustCompile(`Alias for -([A-Za-z0-9_-]+)`)
	defaultPattern := regexp.MustCompile(`default:\s*([^\]]+)`)

	seen := make(map[string]struct{})
	var flags []schemaFlag
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "-") {
			continue
		}

		match := flagNamePattern.FindStringSubmatch(trimmed)
		if len(match) != 3 {
			continue
		}

		name := match[1]
		description := strings.TrimSpace(match[2])
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}

		typeName := "string"
		lowerDesc := strings.ToLower(description)
		if strings.Contains(lowerDesc, "true/false") {
			typeName = "bool"
		}

		aliases := []string{}
		if m := aliasPattern.FindStringSubmatch(description); len(m) == 2 {
			aliases = append(aliases, m[1])
		}

		defaultValue := ""
		if m := defaultPattern.FindStringSubmatch(description); len(m) == 2 {
			defaultValue = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(strings.TrimSpace(m[1]), ")"), "]"))
			defaultValue = strings.TrimPrefix(strings.TrimSuffix(defaultValue, "."), "(default:")
		}

		flags = append(flags, schemaFlag{
			Name:        name,
			Type:        typeName,
			Required:    strings.Contains(description, "(required)"),
			Aliases:     aliases,
			Description: strings.TrimSpace(strings.Join(strings.Fields(description), " ")),
			Default:     defaultValue,
		})
	}

	return flags
}

// Help returns help text for the schema command
func (c *SchemaCommand) Help() string {
	helpText := `
Usage: hcptf schema <command> [args]

  Print a JSON description of command flags.

Examples:

  hcptf schema workspace create
  hcptf schema run apply
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the schema command
func (c *SchemaCommand) Synopsis() string {
	return "Print command schema"
}
