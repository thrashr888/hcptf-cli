package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcptf-cli/command"
	"github.com/hashicorp/hcptf-cli/internal/router"
	"github.com/mitchellh/cli"
)

const (
	// Version is the main version number
	Version = "0.1.0"

	// VersionPrerelease is a pre-release marker for the version
	VersionPrerelease = "dev"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	// Disable log output by default
	log.SetOutput(os.Stderr)

	// Create the meta object for commands
	meta := command.Meta{
		Color: true,
	}

	// Setup UI with color support
	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		InfoColor:  cli.UiColorGreen,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	meta.Ui = ui

	// Translate URL-like args if present
	args := os.Args[1:]
	r := router.NewRouter(nil) // Pass nil client for now - we don't need validation
	translatedArgs, err := r.TranslateArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %s\n", err.Error())
		return 1
	}

	// Create CLI instance
	c := &cli.CLI{
		Name:       "hcptf",
		Version:    GetVersion(),
		Args:       translatedArgs,
		Commands:   command.Commands(&meta),
		HelpFunc:   cli.BasicHelpFunc("hcptf"),
		HelpWriter: os.Stdout,
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

// GetVersion returns the full version string
func GetVersion() string {
	if VersionPrerelease != "" {
		return fmt.Sprintf("%s-%s", Version, VersionPrerelease)
	}
	return Version
}
