package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/config"
	"github.com/hashicorp/hcptf-cli/internal/validate"
	"github.com/hashicorp/hcptf-cli/internal/output"
	"github.com/mitchellh/cli"
)

// Meta contains the meta-options and functionality shared between all commands
type Meta struct {
	// Color enables colorized output
	Color bool

	// OutputWriter is the destination for table and plain-text formatter output.
	OutputWriter io.Writer

	// ErrorWriter is the destination for formatter JSON encoding errors.
	ErrorWriter io.Writer

	// Fields for command output filtering.
	Fields string

	// DryRun skips API calls for mutation commands.
	DryRun bool

	// JSONInput reads API payloads for mutation commands.
	JSONInput string

	// Ui is the CLI user interface
	Ui cli.Ui

	// client is the cached API client
	client *client.Client

	// clientErr is any error that occurred during client initialization
	clientErr error

	// config is the cached CLI configuration
	config *config.Config

	// configErr is any error that occurred during config loading
	configErr error
}

// Client returns the API client, initializing it if necessary
func (m *Meta) Client() (*client.Client, error) {
	if m.client != nil || m.clientErr != nil {
		return m.client, m.clientErr
	}

	cfg, err := m.Config()
	if err != nil {
		m.clientErr = fmt.Errorf("failed to load config: %w", err)
		return nil, m.clientErr
	}

	m.client, m.clientErr = client.New(cfg)
	return m.client, m.clientErr
}

func (m *Meta) formatterWriter() io.Writer {
	if m.OutputWriter != nil {
		return m.OutputWriter
	}
	if mock, ok := m.Ui.(*cli.MockUi); ok && mock.OutputWriter != nil {
		return mock.OutputWriter
	}

	return os.Stdout
}

func (m *Meta) formatterErrorWriter() io.Writer {
	if m.ErrorWriter != nil {
		return m.ErrorWriter
	}
	if mock, ok := m.Ui.(*cli.MockUi); ok && mock.ErrorWriter != nil {
		return mock.ErrorWriter
	}

	return os.Stderr
}

func (m *Meta) NewFormatter(format string) *output.Formatter {
	f := output.NewFormatterWithWriters(format, m.formatterWriter(), m.formatterErrorWriter())
	f.SetFields(parseFields(m.Fields))
	return f
}

// Config returns the CLI configuration, loading it if necessary
func (m *Meta) Config() (*config.Config, error) {
	if m.config != nil || m.configErr != nil {
		return m.config, m.configErr
	}

	m.config, m.configErr = config.Load()
	return m.config, m.configErr
}

// FlagSet returns a FlagSet with common flags
func (m *Meta) FlagSet(name string) *flag.FlagSet {
	f := flag.NewFlagSet(name, flag.ContinueOnError)

	// Common agent-friendly flags
	f.StringVar(&m.Fields, "fields", "", "Output field filter (comma-separated list)")
	f.BoolVar(&m.DryRun, "dry-run", false, "Print planned mutation without API call")
	f.StringVar(&m.JSONInput, "json-input", "", "JSON payload for mutation commands")

	f.Usage = func() {}
	return f
}

// ParseJSONInput parses JSON from a file path, inline JSON, or stdin.
func (m *Meta) ParseJSONInput(target interface{}) error {
	if m.JSONInput == "-" {
		body, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, target)
	}

	if strings.HasPrefix(m.JSONInput, "@") {
		path := strings.TrimPrefix(m.JSONInput, "@")
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, target)
	}

	return json.Unmarshal([]byte(m.JSONInput), target)
}

// ValidateID validates an ID-like value and emits an error message on failure.
func (m *Meta) ValidateID(value, flagName string) bool {
	if err := validate.ID(value, flagName); err != nil {
		m.Ui.Error(fmt.Sprintf("Error: %s", err))
		return false
	}
	return true
}

// ValidateName validates a name-like value and emits an error message on failure.
func (m *Meta) ValidateName(value, flagName string) bool {
	if err := validate.Name(value, flagName); err != nil {
		m.Ui.Error(fmt.Sprintf("Error: %s", err))
		return false
	}
	return true
}

// ValidateString validates a text value and emits an error message on failure.
func (m *Meta) ValidateString(value, flagName string) bool {
	if err := validate.SafeString(value, flagName); err != nil {
		m.Ui.Error(fmt.Sprintf("Error: %s", err))
		return false
	}
	return true
}

func parseFields(raw string) []string {
	var fields []string
	for _, field := range strings.Split(raw, ",") {
		f := strings.TrimSpace(field)
		if f == "" {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

// AutocompleteFlags returns a set of flags for autocomplete
func (m *Meta) AutocompleteFlags() map[string]string {
	return map[string]string{
		"-output": "Output format (table, json)",
	}
}

// ColoredOutput wraps a message with color codes if color is enabled
func (m *Meta) ColoredOutput(color string, message string) string {
	if !m.Color {
		return message
	}
	return fmt.Sprintf("%s%s\033[0m", color, message)
}

// ErrorColor returns the error color code
func (m *Meta) ErrorColor() string {
	return "\033[31m" // Red
}

// SuccessColor returns the success color code
func (m *Meta) SuccessColor() string {
	return "\033[32m" // Green
}

// WarnColor returns the warning color code
func (m *Meta) WarnColor() string {
	return "\033[33m" // Yellow
}

// InfoColor returns the info color code
func (m *Meta) InfoColor() string {
	return "\033[36m" // Cyan
}
