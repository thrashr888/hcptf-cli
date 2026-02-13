package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestMetaFlagSet(t *testing.T) {
	m := Meta{Ui: cli.NewMockUi()}
	fs := m.FlagSet("test")

	if fs == nil {
		t.Fatal("expected non-nil FlagSet")
	}

	fs.String("example", "", "test flag")
	if err := fs.Parse([]string{"-example=hello"}); err != nil {
		t.Fatalf("expected successful parse, got %v", err)
	}
}

func TestMetaAutocompleteFlags(t *testing.T) {
	m := Meta{}
	flags := m.AutocompleteFlags()
	if flags == nil {
		t.Fatal("expected non-nil flags map")
	}
	if _, ok := flags["-output"]; !ok {
		t.Fatal("expected -output in autocomplete flags")
	}
}

func TestMetaColoredOutput(t *testing.T) {
	m := Meta{Color: false}
	out := m.ColoredOutput("\033[31m", "hello")
	if out != "hello" {
		t.Fatalf("expected plain text when color disabled, got %q", out)
	}

	m.Color = true
	out = m.ColoredOutput("\033[31m", "hello")
	if out != "\033[31mhello\033[0m" {
		t.Fatalf("expected colored text, got %q", out)
	}
}

func TestMetaColorMethods(t *testing.T) {
	m := Meta{}

	if c := m.ErrorColor(); c == "" {
		t.Fatal("ErrorColor should not be empty")
	}
	if c := m.SuccessColor(); c == "" {
		t.Fatal("SuccessColor should not be empty")
	}
	if c := m.WarnColor(); c == "" {
		t.Fatal("WarnColor should not be empty")
	}
	if c := m.InfoColor(); c == "" {
		t.Fatal("InfoColor should not be empty")
	}
}

func TestMetaClientCachesResult(t *testing.T) {
	ui := cli.NewMockUi()
	m := newTestMeta(ui)

	c1, err1 := m.Client()
	c2, err2 := m.Client()

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v, %v", err1, err2)
	}
	if c1 != c2 {
		t.Fatal("expected same client instance on second call")
	}
}
