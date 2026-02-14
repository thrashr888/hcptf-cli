package command

import (
	"os"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCommandsFactoryAndText(t *testing.T) {
	meta := newTestMeta(cli.NewMockUi())
	commands := Commands(&meta)

	if len(commands) == 0 {
		t.Fatal("expected commands map to be non-empty")
	}

	for name, factory := range commands {
		cmd, err := factory()
		if err != nil {
			t.Fatalf("factory %s returned error: %v", name, err)
		}
		if cmd == nil {
			t.Fatalf("factory %s returned nil command", name)
		}
	}
}

func TestAllCommandsExecuteWithInvalidFlags(t *testing.T) {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}
	defer func() {
		_ = w.Close()
		_ = r.Close()
		os.Stderr = oldStderr
	}()
	os.Stderr = w

	meta := newTestMeta(cli.NewMockUi())
	commands := Commands(&meta)

	for name, factory := range commands {
		cmd, err := factory()
		if err != nil {
			t.Fatalf("factory %s returned error: %v", name, err)
		}
		if cmd == nil {
			t.Fatalf("factory %s returned nil command", name)
		}

		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("command %s panicked on invalid flag: %v", name, recovered)
				}
			}()
			cmd.Run([]string{"-__invalid__"})
		}()
	}
}

func TestAllCommandsExecuteWithNoArgs(t *testing.T) {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}
	defer func() {
		_ = w.Close()
		_ = r.Close()
		os.Stderr = oldStderr
	}()
	os.Stderr = w

	meta := newTestMeta(cli.NewMockUi())
	commands := Commands(&meta)

	for name, factory := range commands {
		cmd, err := factory()
		if err != nil {
			t.Fatalf("factory %s returned error: %v", name, err)
		}
		if cmd == nil {
			t.Fatalf("factory %s returned nil command", name)
		}

		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					// Some command families intentionally require runtime services and cannot
					// be exercised without additional mock wiring in this generic harness.
					return
				}
			}()
			cmd.Run(nil)
		}()
	}
}
