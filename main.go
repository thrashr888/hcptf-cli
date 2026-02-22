package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/hcptf-cli/command"
	"github.com/hashicorp/hcptf-cli/internal/router"
	"github.com/mitchellh/cli"
)

var (
	// Version is the main version number (injected by GoReleaser via -ldflags)
	Version = "0.5.0"

	// VersionPrerelease is a pre-release marker for the version (injected by GoReleaser)
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
		Color:        true,
		OutputWriter: os.Stdout,
		ErrorWriter:  os.Stderr,
	}
	command.SetVersionProvider(GetVersion)

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

	commands := command.Commands(&meta)
	commandPaths := extractCommandPaths(commands)
	getVerbIndex := buildGetVerbIndex(commands)

	// Translate URL-like args if present
	args := os.Args[1:]
	r := router.NewRouter(nil, commandPaths) // Pass nil client for now - we don't need validation
	translatedArgs, err := r.TranslateArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %s\n", err.Error())
		return 1
	}
	inferredArgs, err := inferImplicitGetVerb(translatedArgs, getVerbIndex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return 1
	}
	normalizedArgs := normalizeDeleteConfirmationFlags(inferredArgs, commands)

	// Create CLI instance
	c := &cli.CLI{
		Name:       "hcptf",
		Version:    GetVersion(),
		Args:       normalizedArgs,
		Commands:   commands,
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

func extractCommandPaths(commands map[string]cli.CommandFactory) []string {
	paths := make([]string, 0, len(commands))
	for key := range commands {
		if strings.TrimSpace(key) == "" {
			continue
		}
		paths = append(paths, key)
	}
	sort.Strings(paths)
	return paths
}

type getVerbAvailability struct {
	hasList bool
	hasRead bool
	hasShow bool
}

func buildGetVerbIndex(commands map[string]cli.CommandFactory) map[string]getVerbAvailability {
	index := make(map[string]getVerbAvailability)
	for key := range commands {
		parts := strings.Split(key, " ")
		if len(parts) < 2 {
			continue
		}

		last := parts[len(parts)-1]
		namespace := strings.Join(parts[:len(parts)-1], " ")

		current := index[namespace]
		switch last {
		case "list":
			current.hasList = true
		case "read":
			current.hasRead = true
		case "show":
			current.hasShow = true
		default:
			continue
		}
		index[namespace] = current
	}

	return index
}

func inferImplicitGetVerb(args []string, getVerbIndex map[string]getVerbAvailability) ([]string, error) {
	if len(args) == 0 || len(getVerbIndex) == 0 {
		return args, nil
	}
	if strings.HasPrefix(args[0], "-") || hasHelpFlag(args) {
		return args, nil
	}

	firstFlag := len(args)
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			firstFlag = i
			break
		}
	}

	tokens := args[:firstFlag]
	if len(tokens) == 0 {
		return args, nil
	}

	matchLen := 0
	var match getVerbAvailability
	for i := len(tokens); i >= 1; i-- {
		candidate := strings.Join(tokens[:i], " ")
		if v, ok := getVerbIndex[candidate]; ok {
			matchLen = i
			match = v
			break
		}
	}
	if matchLen == 0 {
		return args, nil
	}
	namespace := strings.Join(tokens[:matchLen], " ")

	if matchLen < len(tokens) {
		return args, nil
	}

	verb, ambiguous := chooseImplicitGetVerb(args[firstFlag:], match)
	if ambiguous {
		available := availableGetVerbs(match)
		return nil, fmt.Errorf(
			"ambiguous operation for %q; specify one of: %s",
			namespace,
			strings.Join(available, ", "),
		)
	}
	if verb == "" {
		return args, nil
	}

	out := make([]string, 0, len(args)+1)
	out = append(out, args[:matchLen]...)
	out = append(out, verb)
	out = append(out, args[matchLen:]...)
	return out, nil
}

func chooseImplicitGetVerb(flags []string, verbs getVerbAvailability) (string, bool) {
	total := 0
	if verbs.hasList {
		total++
	}
	if verbs.hasRead {
		total++
	}
	if verbs.hasShow {
		total++
	}
	if total == 0 {
		return "", false
	}
	if total == 1 {
		if verbs.hasList {
			return "list", false
		}
		if verbs.hasRead {
			return "read", false
		}
		return "show", false
	}

	if hasIdentitySelector(flags) {
		if verbs.hasRead {
			return "read", false
		}
		if verbs.hasShow {
			return "show", false
		}
	}

	if verbs.hasList && hasCollectionSelector(flags) {
		return "list", false
	}

	return "", true
}

func availableGetVerbs(verbs getVerbAvailability) []string {
	out := make([]string, 0, 3)
	if verbs.hasList {
		out = append(out, "list")
	}
	if verbs.hasShow {
		out = append(out, "show")
	}
	if verbs.hasRead {
		out = append(out, "read")
	}
	return out
}

func isExplicitGetVerb(token string) bool {
	switch token {
	case "list", "read", "show":
		return true
	default:
		return false
	}
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "-help" {
			return true
		}
	}
	return false
}

func hasIdentitySelector(flags []string) bool {
	for i := 0; i < len(flags); i++ {
		flag := flags[i]
		switch {
		case strings.HasPrefix(flag, "-id="), strings.HasPrefix(flag, "--id="):
			return true
		case strings.HasPrefix(flag, "-name="), strings.HasPrefix(flag, "--name="):
			return true
		case flag == "-id" || flag == "--id":
			if i+1 < len(flags) && !strings.HasPrefix(flags[i+1], "-") {
				return true
			}
		case flag == "-name" || flag == "--name":
			if i+1 < len(flags) && !strings.HasPrefix(flags[i+1], "-") {
				return true
			}
		}
	}
	return false
}

func hasCollectionSelector(flags []string) bool {
	for _, flag := range flags {
		if !strings.HasPrefix(flag, "-") {
			continue
		}
		key := flagName(flag)
		switch key {
		case "id", "name", "h", "help", "output", "o":
			continue
		default:
			return true
		}
	}
	return false
}

func flagName(flag string) string {
	trimmed := strings.TrimLeft(flag, "-")
	if trimmed == "" {
		return ""
	}
	if eq := strings.Index(trimmed, "="); eq >= 0 {
		return trimmed[:eq]
	}
	return trimmed
}

func normalizeDeleteConfirmationFlags(args []string, commands map[string]cli.CommandFactory) []string {
	if len(args) == 0 {
		return args
	}
	if strings.HasPrefix(args[0], "-") || hasHelpFlag(args) {
		return args
	}

	firstFlag := len(args)
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			firstFlag = i
			break
		}
	}

	commandTokens := args[:firstFlag]
	if len(commandTokens) == 0 {
		return args
	}

	matchLen := 0
	matchKey := ""
	for i := len(commandTokens); i >= 1; i-- {
		candidate := strings.Join(commandTokens[:i], " ")
		if _, ok := commands[candidate]; ok {
			matchLen = i
			matchKey = candidate
			break
		}
	}
	if matchLen == 0 || !strings.HasSuffix(matchKey, " delete") {
		return args
	}

	out := append([]string(nil), args...)
	for i := firstFlag; i < len(out); i++ {
		if out[i] == "-f" || out[i] == "-y" {
			out[i] = "-force"
		}
	}
	return out
}
