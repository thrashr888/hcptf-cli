package output

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	fn()
	w.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	return string(data)
}

func TestNewFormatterDefaultsToTable(t *testing.T) {
	f := NewFormatter("unsupported")
	if got := f.GetFormat(); got != FormatTable {
		t.Fatalf("expected FormatTable, got %s", got)
	}

	f = NewFormatter("json")
	if got := f.GetFormat(); got != FormatJSON {
		t.Fatalf("expected FormatJSON, got %s", got)
	}
}

func TestTableOutputsJSONWhenRequested(t *testing.T) {
	formatter := NewFormatter("json")
	out := captureStdout(t, func() {
		formatter.Table([]string{"name", "id"}, [][]string{{"hcptf", "1"}})
	})

	var rows []map[string]string
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}
	if rows[0]["name"] != "hcptf" || rows[0]["id"] != "1" {
		t.Fatalf("unexpected row data: %#v", rows[0])
	}
}

func TestJSONFallbackInTableMode(t *testing.T) {
	formatter := NewFormatter("table")
	out := captureStdout(t, func() {
		formatter.JSON(map[string]interface{}{"token": "abc123"})
	})

	if out != "token: abc123\n" {
		t.Fatalf("unexpected table-mode json output: %q", out)
	}
}

func TestListRespectsFormats(t *testing.T) {
	list := []string{"one", "two"}

	tableOut := captureStdout(t, func() {
		NewFormatter("table").List(list)
	})
	if tableOut != "one\ntwo\n" {
		t.Fatalf("unexpected table output: %q", tableOut)
	}

	jsonOut := captureStdout(t, func() {
		NewFormatter("json").List(list)
	})

	var decoded []string
	if err := json.Unmarshal([]byte(jsonOut), &decoded); err != nil {
		t.Fatalf("failed to decode list json: %v", err)
	}
	if len(decoded) != 2 || decoded[0] != "one" || decoded[1] != "two" {
		t.Fatalf("unexpected list data: %v", decoded)
	}
}
