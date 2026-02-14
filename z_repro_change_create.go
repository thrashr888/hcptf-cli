//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/hashicorp/hcptf-cli/command"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/config"
	"github.com/mitchellh/cli"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			fmt.Fprint(w, `{"ok":true}`)
		case "/api/v2/ping?":
			fmt.Fprint(w, `{"ok":true}`)
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			fmt.Fprint(w, `{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`)
		case "/api/v2/organizations/my-org/explorer/bulk-actions":
			fmt.Fprint(w, `{"data":{"id":"ba-123","type":"bulk-actions","attributes":{"organization_id":"my-org","action_type":"change_requests","action_inputs":{"subject":"Fix","message":"Please update"}}}`)
		case "/api/v2/change-requests/cr-1":
			fmt.Fprint(w, `{"data":{"id":"cr-1","type":"change-requests","attributes":{"subject":"Fix","message":"Please update","archived-by":null,"archived-at":null,"created-at":"2024-01-01T00:00:00Z","updated-at":"2024-01-02T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}`)
		default:
			fmt.Printf("unexpected: %s %s\n", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"errors":[{"status":"404"}]}`)
		}
	}))
	defer server.Close()

	parsed, _ := url.Parse(server.URL)
	cfg := &config.Config{
		Credentials: map[string]*config.Credential{
			parsed.Host: &config.Credential{
				Hostname: parsed.Host,
				Token:    "test-token",
			},
		},
	}

	apiClient, err := client.New(cfg)
	if err != nil {
		fmt.Println("new client err", err)
		os.Exit(1)
	}

	ui := cli.NewMockUi()
	cmd := &command.ChangeRequestCreateCommand{Meta: command.Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-subject=Fix", "-message=Please update"})
	fmt.Println("create exit", code)
	fmt.Printf("out=%q\n", ui.OutputWriter.String())
	fmt.Printf("err=%q\n", ui.ErrorWriter.String())

	readUI := cli.NewMockUi()
	readCmd := &command.ChangeRequestReadCommand{Meta: command.Meta{Ui: readUI, client: apiClient}}
	code2 := readCmd.Run([]string{"-id=cr-1"})
	fmt.Println("read exit", code2)
	fmt.Printf("read_out=%q\n", readUI.OutputWriter.String())
	fmt.Printf("read_err=%q\n", readUI.ErrorWriter.String())
}
