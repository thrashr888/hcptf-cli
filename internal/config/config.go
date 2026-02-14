package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Config represents the CLI configuration
type Config struct {
	// Credentials is a map of hostname to credentials
	Credentials map[string]*Credential `hcl:"credentials,block"`

	// DefaultOrganization is the default organization to use
	DefaultOrganization string `hcl:"default_organization,optional"`

	// OutputFormat is the default output format (table, json)
	OutputFormat string `hcl:"output_format,optional"`
}

type fileConfig struct {
	Credentials         []*Credential `hcl:"credentials,block"`
	DefaultOrganization string        `hcl:"default_organization,optional"`
	OutputFormat        string        `hcl:"output_format,optional"`
}

// Credential represents credentials for a specific Terraform instance
type Credential struct {
	Hostname string `hcl:"hostname,label"`
	Token    string `hcl:"token"`
}

// TerraformCredentials represents the Terraform CLI credentials file format
type TerraformCredentials struct {
	Credentials map[string]TerraformCredential `json:"credentials"`
}

// TerraformCredential represents a single credential in the Terraform CLI format
type TerraformCredential struct {
	Token string `json:"token"`
}

// Load loads the configuration from the default location or environment
func Load() (*Config, error) {
	configPath := GetConfigPath()

	var config Config
	config.Credentials = make(map[string]*Credential)
	config.OutputFormat = "table"

	// Try to load from hcptfrc first
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}

		virtualPath := configPath
		switch filepath.Ext(virtualPath) {
		case ".hcl", ".json":
		default:
			virtualPath += ".hcl"
		}

		var diskConfig fileConfig
		if err := hclsimple.Decode(virtualPath, data, nil, &diskConfig); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}

		for _, cred := range diskConfig.Credentials {
			config.Credentials[cred.Hostname] = cred
		}

		if diskConfig.DefaultOrganization != "" {
			config.DefaultOrganization = diskConfig.DefaultOrganization
		}

		if diskConfig.OutputFormat != "" {
			config.OutputFormat = diskConfig.OutputFormat
		}
	}

	// Also try to load credentials from Terraform CLI credentials file
	tfCredsPath := GetTerraformCredentialsPath()
	if _, err := os.Stat(tfCredsPath); err == nil {
		tfCreds, err := loadTerraformCredentials(tfCredsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load Terraform credentials from %q: %w", tfCredsPath, err)
		}

		// Merge Terraform credentials into config
		// hcptfrc credentials take precedence
		for hostname, cred := range tfCreds.Credentials {
			if _, exists := config.Credentials[hostname]; !exists {
				config.Credentials[hostname] = &Credential{
					Hostname: hostname,
					Token:    cred.Token,
				}
			}
		}
	}

	// Set defaults
	if config.OutputFormat == "" {
		config.OutputFormat = "table"
	}

	if config.Credentials == nil {
		config.Credentials = make(map[string]*Credential)
	}

	return &config, nil
}

// GetToken returns the authentication token for the given hostname
// Priority: TFE_TOKEN env var > HCPTF_TOKEN env var > config files
func (c *Config) GetToken(hostname string) string {
	// Check TFE_TOKEN environment variable first (standard Terraform CLI env var)
	if token := os.Getenv("TFE_TOKEN"); token != "" {
		return token
	}

	// Check HCPTF_TOKEN environment variable (hcptf-specific)
	if token := os.Getenv("HCPTF_TOKEN"); token != "" {
		return token
	}

	// Check config files (both hcptfrc and Terraform CLI credentials)
	if cred, ok := c.Credentials[hostname]; ok {
		return cred.Token
	}

	return ""
}

// GetAddress returns the API address to use
// Priority: HCPTF_ADDRESS env var > default
func GetAddress() string {
	// Check HCPTF_ADDRESS first (new standard)
	if addr := os.Getenv("HCPTF_ADDRESS"); addr != "" {
		return addr
	}
	// Fall back to TFE_ADDRESS for compatibility
	if addr := os.Getenv("TFE_ADDRESS"); addr != "" {
		return addr
	}
	return "https://app.terraform.io"
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	if path := os.Getenv("HCPTF_CONFIG"); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".hcptfrc")
}

// GetTerraformCredentialsPath returns the path to the Terraform CLI credentials file
func GetTerraformCredentialsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".terraform.d", "credentials.tfrc.json")
}

// loadTerraformCredentials loads credentials from the Terraform CLI credentials file
func loadTerraformCredentials(path string) (*TerraformCredentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Terraform credentials: %w", err)
	}

	var creds TerraformCredentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Terraform credentials: %w", err)
	}

	return &creds, nil
}

// LoadTerraformCredentialsFile loads the Terraform CLI credentials file
// This is a public function for use by login/logout commands
func LoadTerraformCredentialsFile() (*TerraformCredentials, error) {
	path := GetTerraformCredentialsPath()
	return loadTerraformCredentials(path)
}

// ValidateToken validates a token by making a test API call
func ValidateToken(hostname, token string) error {
	address := fmt.Sprintf("https://%s", hostname)

	config := &tfe.Config{
		Address:    address,
		Token:      token,
		HTTPClient: http.DefaultClient,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Make a simple API call to validate the token
	// We'll try to read the account details
	_, err = client.Users.ReadCurrent(nil)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	return nil
}

// SaveCredential saves a credential to the Terraform CLI credentials file
func SaveCredential(hostname, token string) error {
	credsPath := GetTerraformCredentialsPath()

	// Create directory if it doesn't exist
	dir := filepath.Dir(credsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	// Load existing credentials or create new
	creds, err := loadTerraformCredentials(credsPath)
	if err != nil {
		// If file doesn't exist, create new credentials
		creds = &TerraformCredentials{
			Credentials: make(map[string]TerraformCredential),
		}
	}

	// Add or update the credential
	creds.Credentials[hostname] = TerraformCredential{
		Token: token,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(credsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// RemoveCredential removes a credential from the Terraform CLI credentials file
func RemoveCredential(hostname string) error {
	credsPath := GetTerraformCredentialsPath()

	// Load existing credentials
	creds, err := loadTerraformCredentials(credsPath)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	// Remove the credential
	delete(creds.Credentials, hostname)

	// If no credentials left, remove the file
	if len(creds.Credentials) == 0 {
		if err := os.Remove(credsPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove credentials file: %w", err)
		}
		return nil
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(credsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}
