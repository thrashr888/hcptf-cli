package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
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
	format Format
	out    io.Writer
	err    io.Writer
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
		format: f,
		out:    out,
		err:    err,
	}
}

// Table outputs data in table format
func (f *Formatter) Table(headers []string, rows [][]string) {
	if f.format == FormatJSON {
		// Convert table to JSON
		var data []map[string]string
		for _, row := range rows {
			item := make(map[string]string)
			for i, header := range headers {
				if i < len(row) {
					item[header] = row[i]
				}
			}
			data = append(data, item)
		}
		f.JSON(data)
		return
	}

	// Table format with tablewriter
	table := tablewriter.NewWriter(f.out)
	table.Header(headers)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}

// TableWithFullRows outputs data in table format with truncated display values,
// but uses full (untruncated) values for JSON output.
func (f *Formatter) TableWithFullRows(headers []string, displayRows [][]string, fullRows [][]string) {
	if f.format == FormatJSON {
		var data []map[string]string
		for _, row := range fullRows {
			item := make(map[string]string)
			for i, header := range headers {
				if i < len(row) {
					item[header] = row[i]
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
		f.JSON(data)
		return
	}

	// Table format - print as aligned key-value pairs
	var maxKeyLen int
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		padding := strings.Repeat(" ", maxKeyLen-len(k))
		fmt.Fprintf(f.out, "%s:%s %v\n", k, padding, data[k])
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

// GetFormat returns the current output format
func (f *Formatter) GetFormat() Format {
	return f.format
}
