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

func TestKeyValueTable(t *testing.T) {
	formatter := NewFormatter("table")
	out := captureStdout(t, func() {
		formatter.KeyValue(map[string]interface{}{
			"ID":   "ws-123",
			"Name": "prod",
		})
	})

	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !contains(out, "ID") || !contains(out, "ws-123") {
		t.Fatalf("expected key-value pairs in output, got %q", out)
	}
}

func TestKeyValueJSON(t *testing.T) {
	formatter := NewFormatter("json")
	out := captureStdout(t, func() {
		formatter.KeyValue(map[string]interface{}{
			"ID":   "ws-123",
			"Name": "prod",
		})
	})

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", out, err)
	}
	if data["ID"] != "ws-123" {
		t.Fatalf("expected ID ws-123, got %v", data["ID"])
	}
}

func TestTableRendersRows(t *testing.T) {
	formatter := NewFormatter("table")
	out := captureStdout(t, func() {
		formatter.Table([]string{"Name", "ID"}, [][]string{
			{"prod", "ws-1"},
			{"staging", "ws-2"},
		})
	})

	if out == "" {
		t.Fatal("expected non-empty table output")
	}
	if !contains(out, "prod") || !contains(out, "staging") {
		t.Fatalf("expected row data in output, got %q", out)
	}
}

func TestJSONEncodesDirect(t *testing.T) {
	formatter := NewFormatter("json")
	out := captureStdout(t, func() {
		formatter.JSON([]string{"a", "b"})
	})

	var decoded []string
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("expected valid JSON array, got %q: %v", out, err)
	}
	if len(decoded) != 2 || decoded[0] != "a" {
		t.Fatalf("unexpected data: %v", decoded)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && stringContains(s, substr)
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
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
