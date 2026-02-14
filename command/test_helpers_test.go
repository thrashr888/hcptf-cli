package command

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/mitchellh/cli"
)

var (
	mockOutputWriterMu sync.Mutex
	lastMockOutput     io.Writer
)

func newTestMeta(ui cli.Ui) Meta {
	outputWriter, errorWriter := testUICaptureWriters(ui)
	mockOutputWriterMu.Lock()
	lastMockOutput = outputWriter
	mockOutputWriterMu.Unlock()

	return Meta{
		Ui:           ui,
		client:       &client.Client{},
		OutputWriter: outputWriter,
		ErrorWriter:  errorWriter,
	}
}

func testUICaptureWriters(ui cli.Ui) (io.Writer, io.Writer) {
	if mock, ok := ui.(*cli.MockUi); ok {
		return mock.OutputWriter, mock.ErrorWriter
	}

	return os.Stdout, os.Stderr
}

// testMeta is an alias for newTestMeta for backwards compatibility
func testMeta(t *testing.T, ui cli.Ui) Meta {
	t.Helper()
	return newTestMeta(ui)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func captureStdout(t *testing.T, fn func() int) (string, int) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	mockOutputWriterMu.Lock()
	capturedMockOutput := lastMockOutput
	mockOutputWriterMu.Unlock()
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()
	code := fn()
	w.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}
	output := string(data)
	if output == "" && capturedMockOutput != nil {
		output = captureWriterOutput(capturedMockOutput)
	}

	return output, code
}

type stringer interface {
	String() string
}

func captureWriterOutput(writer io.Writer) string {
	if writer == nil {
		return ""
	}

	if buffer, ok := writer.(*bytes.Buffer); ok {
		return buffer.String()
	}
	if stringerWriter, ok := writer.(stringer); ok {
		return stringerWriter.String()
	}

	return ""
}
