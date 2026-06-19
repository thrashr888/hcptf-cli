package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const EnvFileVariable = "HCPTF_ENV_FILE"

// LoadDotEnv loads connection settings from an explicit env file or a default
// .env in the current working directory. Existing environment variables win.
func LoadDotEnv() error {
	if path := os.Getenv(EnvFileVariable); path != "" {
		return LoadDotEnvFile(path)
	}

	if _, err := os.Stat(".env"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to stat .env: %w", err)
	}

	return LoadDotEnvFile(".env")
}

// LoadDotEnvFile loads variables from path without overriding exported env vars.
func LoadDotEnvFile(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load env file %q: %w", path, err)
	}
	return nil
}
