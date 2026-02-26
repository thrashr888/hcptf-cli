package output

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestNewFormatterWithWriters(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("table", out, errOut)

	formatter.List([]string{"alpha", "beta"})
	if out.String() != "alpha\nbeta\n" {
		t.Fatalf("expected table output, got %q", out.String())
	}
	if errOut.Len() != 0 {
		t.Fatalf("expected no errors, got %q", errOut.String())
	}
}

func TestNewFormatterWithWritersDefaultsForNilWriters(t *testing.T) {
	out := captureStdout(t, func() {
		formatter := NewFormatterWithWriters("table", nil, nil)
		formatter.List([]string{"alpha", "beta"})
	})
	if out != "alpha\nbeta\n" {
		t.Fatalf("expected stdout fallback output, got %q", out)
	}
}

type errorWriter struct {
	buffer      bytes.Buffer
	errToReturn error
}

func (w *errorWriter) Write(p []byte) (int, error) {
	w.buffer.Write(p)
	return len(p), w.errToReturn
}

func TestJSONEncodingErrorWritesToErrorWriter(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &errorWriter{errToReturn: errors.New("forced write error")}
	formatter := NewFormatterWithWriters("json", out, errOut)

	formatter.JSON(func() string { return "unsupported type" })

	if out.Len() != 0 {
		t.Fatalf("expected no normal JSON output on encode failure, got %q", out.String())
	}
	if errOut.buffer.Len() == 0 {
		t.Fatal("expected error output to be written")
	}
}

func TestTableOutputsJSONWhenRequested(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("json", out, &bytes.Buffer{})
	captureStdout(t, func() {
		formatter.Table([]string{"name", "id"}, [][]string{{"hcptf", "1"}})
	})

	var rows []map[string]string
	if err := json.Unmarshal([]byte(out.String()), &rows); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}
	if rows[0]["name"] != "hcptf" || rows[0]["id"] != "1" {
		t.Fatalf("unexpected row data: %#v", rows[0])
	}
}

func TestJSONFallbackInTableMode(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("table", out, &bytes.Buffer{})
	captureStdout(t, func() {
		formatter.JSON(map[string]interface{}{"token": "abc123"})
	})

	if out.String() != "token: abc123\n" {
		t.Fatalf("unexpected table-mode json output: %q", out.String())
	}
}

func TestKeyValueTable(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("table", out, &bytes.Buffer{})
	captureStdout(t, func() {
		formatter.KeyValue(map[string]interface{}{
			"ID":   "ws-123",
			"Name": "prod",
		})
	})

	if out.String() == "" {
		t.Fatal("expected non-empty output")
	}
	if !contains(out.String(), "ID") || !contains(out.String(), "ws-123") {
		t.Fatalf("expected key-value pairs in output, got %q", out.String())
	}
}

func TestKeyValueJSON(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("json", out, &bytes.Buffer{})
	captureStdout(t, func() {
		formatter.KeyValue(map[string]interface{}{
			"ID":   "ws-123",
			"Name": "prod",
		})
	})

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out.String()), &data); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", out.String(), err)
	}
	if data["ID"] != "ws-123" {
		t.Fatalf("expected ID ws-123, got %v", data["ID"])
	}
}

func TestTableRendersRows(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("table", out, &bytes.Buffer{})
	captureStdout(t, func() {
		formatter.Table([]string{"Name", "ID"}, [][]string{
			{"prod", "ws-1"},
			{"staging", "ws-2"},
		})
	})

	if out.String() == "" {
		t.Fatal("expected non-empty table output")
	}
	if !contains(out.String(), "prod") || !contains(out.String(), "staging") {
		t.Fatalf("expected row data in output, got %q", out.String())
	}
}

func TestJSONEncodesDirect(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("json", out, &bytes.Buffer{})
	captureStdout(t, func() {
		formatter.JSON([]string{"a", "b"})
	})

	var decoded []string
	if err := json.Unmarshal([]byte(out.String()), &decoded); err != nil {
		t.Fatalf("expected valid JSON array, got %q: %v", out.String(), err)
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

func TestTableWithFullRowsUsesFullValuesForJSON(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("json", out, &bytes.Buffer{})

	headers := []string{"Key", "Value"}
	displayRows := [][]string{{"kafka_config", `{"broker_count...`}}
	fullRows := [][]string{{"kafka_config", `{"broker_count_per_az":1,"version":"3.5.1"}`}}

	formatter.TableWithFullRows(headers, displayRows, fullRows)

	var rows []map[string]string
	if err := json.Unmarshal(out.Bytes(), &rows); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}
	if rows[0]["Value"] != `{"broker_count_per_az":1,"version":"3.5.1"}` {
		t.Fatalf("expected full value in JSON, got %q", rows[0]["Value"])
	}
}

func TestTableWithFullRowsUsesDisplayRowsForTable(t *testing.T) {
	out := &bytes.Buffer{}
	formatter := NewFormatterWithWriters("table", out, &bytes.Buffer{})

	headers := []string{"Key", "Value"}
	displayRows := [][]string{{"kafka_config", `{"broker_count...`}}
	fullRows := [][]string{{"kafka_config", `{"broker_count_per_az":1,"version":"3.5.1"}`}}

	formatter.TableWithFullRows(headers, displayRows, fullRows)

	if !contains(out.String(), `{"broker_count...`) {
		t.Fatalf("expected truncated value in table, got %q", out.String())
	}
}

func TestListRespectsFormats(t *testing.T) {
	list := []string{"one", "two"}
	outBuf := &bytes.Buffer{}

	tableFormatter := NewFormatterWithWriters("table", outBuf, &bytes.Buffer{})
	tableFormatter.List(list)
	if outBuf.String() != "one\ntwo\n" {
		t.Fatalf("unexpected table output: %q", outBuf.String())
	}

	outBuf.Reset()
	jsonFormatter := NewFormatterWithWriters("json", outBuf, &bytes.Buffer{})
	jsonFormatter.List(list)

	var decoded []string
	if err := json.Unmarshal([]byte(outBuf.String()), &decoded); err != nil {
		t.Fatalf("failed to decode list json: %v", err)
	}
	if len(decoded) != 2 || decoded[0] != "one" || decoded[1] != "two" {
		t.Fatalf("unexpected list data: %v", decoded)
	}
}
