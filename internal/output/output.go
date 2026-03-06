package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// Format represents the output format
type Format string

const (
	// FormatTable outputs data in table format
	FormatTable Format = "table"

	// FormatJSON outputs data in JSON format
	FormatJSON Format = "json"
)

// Formatter handles output formatting
type Formatter struct {
	format     Format
	out        io.Writer
	err        io.Writer
	fields     []string
	fieldIndex map[string]int
}

// NewFormatter creates a new output formatter
func NewFormatter(format string) *Formatter {
	return NewFormatterWithWriters(format, os.Stdout, os.Stderr)
}

// NewFormatterWithWriters creates a new output formatter with explicit output streams.
func NewFormatterWithWriters(format string, out io.Writer, err io.Writer) *Formatter {
	if out == nil {
		out = os.Stdout
	}

	if err == nil {
		err = os.Stderr
	}

	f := Format(format)
	if f != FormatTable && f != FormatJSON {
		f = FormatTable // Default to table
	}

	return &Formatter{
		format:     f,
		out:        out,
		err:        err,
		fieldIndex: map[string]int{},
	}
}

// SetFields sets output field filtering by key/header name.
func (f *Formatter) SetFields(fields []string) {
	f.fields = fields
	f.fieldIndex = make(map[string]int, len(fields))
	for i, field := range fields {
		f.fieldIndex[field] = i
	}
}

func (f *Formatter) selectedHeaders(headers []string) ([]int, []string) {
	if len(f.fields) == 0 {
		indexes := make([]int, len(headers))
		for i := range headers {
			indexes[i] = i
		}
		return indexes, headers
	}

	indexes := make([]int, 0, len(headers))
	selected := make([]string, 0, len(headers))
	for i, h := range headers {
		if _, ok := f.fieldIndex[h]; ok {
			indexes = append(indexes, i)
			selected = append(selected, h)
		}
	}
	return indexes, selected
}

// filterRow keeps only columns in selected indexes.
func (f *Formatter) filterRow(row []string, indexes []int) []string {
	result := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		if idx < len(row) {
			result = append(result, row[idx])
		} else {
			result = append(result, "")
		}
	}
	return result
}

// filterMap keeps only requested keys.
func (f *Formatter) filterMap(data map[string]interface{}) map[string]interface{} {
	if len(f.fields) == 0 {
		return data
	}

	filtered := make(map[string]interface{}, len(data))
	for _, field := range f.fields {
		if val, ok := data[field]; ok {
			filtered[field] = val
		}
	}
	return filtered
}

// Table outputs data in table format
func (f *Formatter) Table(headers []string, rows [][]string) {
	if f.format == FormatJSON {
		// Convert table to JSON
		filteredHeadersIdx, filteredHeaders := f.selectedHeaders(headers)
		var data []map[string]string
		for _, row := range rows {
			item := make(map[string]string)
			for j, i := range filteredHeadersIdx {
				if i < len(row) {
					item[filteredHeaders[j]] = row[i]
				}
			}
			data = append(data, item)
		}
		f.JSON(data)
		return
	}

	indexes, filteredHeaders := f.selectedHeaders(headers)
	if len(filteredHeaders) == 0 {
		return
	}

	table := tablewriter.NewTable(f.out, tablewriter.WithHeaderAutoFormat(tw.Off))
	table.Header(filteredHeaders)
	for _, row := range rows {
		table.Append(f.filterRow(row, indexes))
	}
	table.Render()
}

// TableWithFullRows outputs data in table format with truncated display values,
// but uses full (untruncated) values for JSON output.
func (f *Formatter) TableWithFullRows(headers []string, displayRows [][]string, fullRows [][]string) {
	if f.format == FormatJSON {
		filteredHeadersIdx, filteredHeaders := f.selectedHeaders(headers)
		var data []map[string]string
		for _, row := range fullRows {
			item := make(map[string]string)
			for j, i := range filteredHeadersIdx {
				if i < len(row) {
					item[filteredHeaders[j]] = row[i]
				}
			}
			data = append(data, item)
		}
		f.JSON(data)
		return
	}

	f.Table(headers, displayRows)
}

// JSON outputs data in JSON format
func (f *Formatter) JSON(data interface{}) {
	if f.format == FormatTable {
		// If table format requested but JSON called, print simple key-value
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				fmt.Fprintf(f.out, "%s: %v\n", k, v)
			}
			return
		}
	}

	encoder := json.NewEncoder(f.out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(f.err, "Error encoding JSON: %v\n", err)
	}
}

// KeyValue outputs key-value pairs
func (f *Formatter) KeyValue(data map[string]interface{}) {
	if f.format == FormatJSON {
		f.JSON(f.filterMap(data))
		return
	}

	filtered := f.filterMap(data)
	var keys []string
	if len(f.fields) > 0 {
		for _, field := range f.fields {
			if _, ok := filtered[field]; ok {
				keys = append(keys, field)
			}
		}
	} else {
		keys = make([]string, 0, len(filtered))
		for k := range filtered {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	var maxKeyLen int
	for _, k := range keys {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for _, k := range keys {
		padding := strings.Repeat(" ", maxKeyLen-len(k))
		fmt.Fprintf(f.out, "%s:%s %s\n", k, padding, formatValue(filtered[k]))
	}
}

// List outputs a simple list of strings
func (f *Formatter) List(items []string) {
	if f.format == FormatJSON {
		f.JSON(items)
		return
	}

	for _, item := range items {
		fmt.Fprintln(f.out, item)
	}
}

// formatValue converts a value to a human-readable string for table output.
// Struct and pointer-to-struct values are JSON-encoded instead of using Go's
// default %v formatting (which produces unreadable &{...} output).
func formatValue(v interface{}) string {
	if v == nil {
		return "<nil>"
	}

	rv := reflect.ValueOf(v)

	// Dereference pointers
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "<nil>"
		}
		rv = rv.Elem()
	}

	// JSON-encode structs and maps for readable output
	if rv.Kind() == reflect.Struct || rv.Kind() == reflect.Map {
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}

	return fmt.Sprintf("%v", v)
}

// GetFormat returns the current output format
func (f *Formatter) GetFormat() Format {
	return f.format
}
