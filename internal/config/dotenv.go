package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const EnvFileVariable = "HCPTF_ENV_FILE"
const DisableEnvFileVariable = "HCPTF_DISABLE_ENV_FILE"

// LoadDotEnv loads connection settings from an explicit env file or a default
// .env in the current working directory, ancestor directories, or user-level
// hcptf env files. Existing environment variables win.
func LoadDotEnv() error {
	if os.Getenv(DisableEnvFileVariable) != "" {
		return nil
	}

	if path := os.Getenv(EnvFileVariable); path != "" {
		return LoadDotEnvFile(path)
	}

	for _, path := range DotEnvCandidatePaths() {
		info, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("failed to stat env file %q: %w", path, err)
		}
		if info.IsDir() {
			continue
		}
		if err := LoadDotEnvFile(path); err != nil {
			return err
		}
	}

	return nil
}

// LoadDotEnvFile loads variables from path without overriding exported env vars.
func LoadDotEnvFile(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load env file %q: %w", path, err)
	}
	return nil
}

// DotEnvCandidatePaths returns env files from highest to lowest precedence.
// godotenv.Load preserves existing variables, so loading in this order lets
// project-specific files win while still allowing user-level defaults.
func DotEnvCandidatePaths() []string {
	seen := map[string]struct{}{}
	var paths []string
	add := func(path string) {
		if path == "" {
			return
		}
		cleaned := filepath.Clean(path)
		if _, ok := seen[cleaned]; ok {
			return
		}
		seen[cleaned] = struct{}{}
		paths = append(paths, cleaned)
	}

	if cwd, err := os.Getwd(); err == nil {
		home, _ := os.UserHomeDir()
		for {
			add(filepath.Join(cwd, ".env"))
			if cwd == home || cwd == filepath.Dir(cwd) {
				break
			}
			cwd = filepath.Dir(cwd)
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		add(filepath.Join(home, ".hcptf.env"))
		add(filepath.Join(home, ".config", "hcptf", "env"))
	}

	return paths
}
