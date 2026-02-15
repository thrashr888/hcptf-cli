package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

func executeAPIRequest(client *client.Client, method, endpoint string, body io.Reader) ([]byte, int, error) {
	baseURL := strings.TrimSuffix(client.GetAddress(), "/")
	path := strings.TrimSpace(endpoint)
	fullURL := baseURL + path

	req, err := http.NewRequestWithContext(client.Context(), method, fullURL, body)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.Token())
	req.Header.Set("Content-Type", "application/vnd.api+json")

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error reading response: %w", err)
	}

	return payload, resp.StatusCode, nil
}

func parseAPIResponse(body []byte) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func printAPIResponse(formatter *output.Formatter, response interface{}) {
	if formatter.GetFormat() == output.FormatJSON || response == nil {
		formatter.JSON(response)
		return
	}

	switch payload := response.(type) {
	case map[string]interface{}:
		if data, ok := payload["data"]; ok {
			if headers, rows, ok := apiDataRows(data); ok {
				formatter.Table(headers, rows)
				return
			}
		}
		formatter.KeyValue(payload)
	case []interface{}:
		if headers, rows, ok := apiDataRows(response); ok {
			formatter.Table(headers, rows)
			return
		}
		formatter.KeyValue(map[string]interface{}{"count": len(payload), "response": payload})
	default:
		formatter.KeyValue(map[string]interface{}{"response": payload})
	}
}

func apiDataRows(data interface{}) ([]string, [][]string, bool) {
	items, ok := data.([]interface{})
	if !ok {
		return nil, nil, false
	}

	if len(items) == 0 {
		return []string{"Message"}, [][]string{{"No data returned"}}, true
	}

	firstItem, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, nil, false
	}

	headers := []string{"ID", "Type"}
	attributes := extractSortedAttributes(firstItem)
	for _, attr := range attributes {
		headers = append(headers, attr)
	}

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id := toPrintable(obj["id"])
		typeName := toPrintable(obj["type"])
		row := []string{id, typeName}
		attrs := map[string]interface{}{}
		if rawAttrs, ok := obj["attributes"].(map[string]interface{}); ok {
			attrs = rawAttrs
		}

		for _, attr := range attributes {
			row = append(row, toPrintable(attrs[attr]))
		}
		rows = append(rows, row)
	}

	return headers, rows, true
}

func extractSortedAttributes(item map[string]interface{}) []string {
	attributes, ok := item["attributes"].(map[string]interface{})
	if !ok {
		return nil
	}

	tokens := make([]string, 0, len(attributes))
	for key := range attributes {
		tokens = append(tokens, key)
	}
	sort.Strings(tokens)
	return tokens
}

func toPrintable(v interface{}) string {
	if v == nil {
		return "-"
	}
	return strings.TrimSpace(strings.Trim(fmt.Sprintf("%v", v), "\n"))
}
