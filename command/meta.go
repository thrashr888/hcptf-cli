package command

import (
	"flag"
	"fmt"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/config"
	"github.com/mitchellh/cli"
)

// Meta contains the meta-options and functionality shared between all commands
type Meta struct {
	// Color enables colorized output
	Color bool

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
	f.Usage = func() {}
	return f
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
