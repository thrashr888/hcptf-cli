# Testing Guide

This document describes the testing standards and practices for the hcptf CLI project, following HashiCorp conventions from [go-tfe](https://github.com/hashicorp/go-tfe).

## Testing Philosophy

- **Comprehensive Coverage**: Aim for high test coverage across all packages
- **Table-Driven Tests**: Use table-driven tests for testing multiple scenarios
- **Mock-Based Testing**: Use mocks for external dependencies (API clients, services)
- **Clear Test Names**: Test names should describe what they test and what behavior is expected
- **Isolation**: Tests should be isolated and not depend on external state

## Test Organization

### File Structure

Tests are colocated with implementation files using the `*_test.go` suffix:
```
command/
├── workspace_list.go
├── workspace_list_test.go
├── workspace_create.go
├── workspace_create_test.go
...
```

### Test Helpers

Common test utilities are centralized in:
- `command/test_helpers_test.go` - Helper functions for command testing
- `command/mocks_test.go` - Mock service implementations

## Writing Command Tests

### Standard Test Pattern

Commands use dependency injection for testability. Here's the standard pattern:

```go
package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

// 1. Define a mock service interface
type mockWorkspaceService struct {
	response    *tfe.WorkspaceList
	err         error
	lastOrg     string
	lastOptions *tfe.WorkspaceListOptions
}

// 2. Implement the service interface
func (m *mockWorkspaceService) List(_ context.Context, organization string, options *tfe.WorkspaceListOptions) (*tfe.WorkspaceList, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

// 3. Create a test constructor with dependency injection
func newWorkspaceListCommand(ui cli.Ui, svc workspaceLister) *WorkspaceListCommand {
	return &WorkspaceListCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
}

// 4. Write test cases
func TestWorkspaceListCommandRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceListCommand(ui, &mockWorkspaceService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}
```

### Required Test Cases

Every command should have tests covering:

1. **Flag Validation**
   - Required flags are enforced
   - Invalid flag values are rejected

2. **API Error Handling**
   - Graceful handling of API errors
   - Appropriate error messages to user

3. **Output Formats**
   - Table output (default)
   - JSON output (-output=json)
   - Empty results handling

4. **Success Cases**
   - Proper service calls with correct parameters
   - Correct output formatting

### Test Naming Convention

Test names follow the pattern: `Test<CommandName><Scenario>`

Examples:
- `TestWorkspaceListRequiresOrganization`
- `TestWorkspaceListHandlesAPIError`
- `TestWorkspaceListOutputsJSON`
- `TestWorkspaceCreateValidatesFlags`

## Table-Driven Tests

For testing multiple scenarios, use table-driven tests:

```go
func TestClient_GetAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{
			name:    "returns default address",
			address: "https://app.terraform.io",
		},
		{
			name:    "returns custom address",
			address: "https://tfe.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				address: tt.address,
			}

			result := client.GetAddress()

			if result != tt.address {
				t.Errorf("expected address '%s', got '%s'", tt.address, result)
			}
		})
	}
}
```

## Mock Services

### Creating Mocks

Mocks implement service interfaces and track calls for verification:

```go
type mockWorkspaceService struct {
	// Return values
	response    *tfe.WorkspaceList
	err         error

	// Call tracking
	lastOrg     string
	lastOptions *tfe.WorkspaceListOptions
	callCount   int
}

func (m *mockWorkspaceService) List(_ context.Context, organization string, options *tfe.WorkspaceListOptions) (*tfe.WorkspaceList, error) {
	m.callCount++
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}
```

### Verifying Mock Calls

```go
func TestWorkspaceListCallsService(t *testing.T) {
	svc := &mockWorkspaceService{
		response: &tfe.WorkspaceList{Items: []*tfe.Workspace{}},
	}
	cmd := newWorkspaceListCommand(cli.NewMockUi(), svc)

	cmd.Run([]string{"-organization=my-org"})

	// Verify the service was called correctly
	if svc.lastOrg != "my-org" {
		t.Errorf("expected org 'my-org', got '%s'", svc.lastOrg)
	}

	if svc.callCount != 1 {
		t.Errorf("expected 1 call, got %d", svc.callCount)
	}
}
```

## Testing Output

### Capturing Stdout

Use the `captureStdout` helper to test output:

```go
func TestWorkspaceListOutput(t *testing.T) {
	ui := cli.NewMockUi()
	// ... setup command ...

	cmd.Run([]string{"-organization=my-org"})

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "expected text") {
		t.Errorf("expected output to contain 'expected text', got: %s", output)
	}
}
```

### Testing JSON Output

```go
func TestWorkspaceListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceService{
		response: &tfe.WorkspaceList{Items: []*tfe.Workspace{{
			ID:   "ws-123",
			Name: "prod",
		}}},
	}
	cmd := newWorkspaceListCommand(ui, svc)

	cmd.Run([]string{"-organization=my-org", "-output=json"})

	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(ui.OutputWriter.String()), &result); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(result))
	}

	if result[0]["Name"] != "prod" {
		t.Errorf("expected name 'prod', got '%v'", result[0]["Name"])
	}
}
```

## Testing Configuration

### Using Environment Variables

```go
func TestConfigFromEnv(t *testing.T) {
	t.Setenv("HCPTF_TOKEN", "test-token")
	t.Setenv("HCPTF_ADDRESS", "https://app.terraform.io")

	// Test code that uses these env vars
}
```

### Using Temp Directories

```go
func TestConfigFile(t *testing.T) {
	tmpDir := t.TempDir() // Automatically cleaned up

	configPath := filepath.Join(tmpDir, "config")
	os.WriteFile(configPath, []byte("..."), 0644)

	// Test code that uses the config file
}
```

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests for Specific Package
```bash
go test ./command/
go test ./internal/client/
```

### Run Tests with Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Test
```bash
go test -run TestWorkspaceList ./command/
```

### Run Tests Verbosely
```bash
go test -v ./...
```

### Run Tests with Race Detector
```bash
go test -race ./...
```

## Integration Tests

Integration tests (marked with build tag `integration`) require real API access:

```go
//go:build integration
// +build integration

package command

func TestWorkspaceListIntegration(t *testing.T) {
	// Requires HCPTF_TOKEN and real API access
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Test against real API
}
```

Run integration tests:
```bash
go test -tags=integration ./...
```

## Test Coverage Goals

- **Internal packages**: 80%+ coverage
- **Command packages**: 70%+ coverage (focus on critical paths)
- **Overall**: 60%+ coverage

Current coverage can be checked with:
```bash
go test -cover ./... | grep -E '^(ok|FAIL)'
```

## Best Practices

1. **Use t.Helper()** in test helpers to improve error reporting
2. **Use t.Parallel()** for independent tests to speed up test runs
3. **Clean up resources** using `t.Cleanup()` or `defer`
4. **Test error paths** not just happy paths
5. **Use meaningful test data** that reflects real-world usage
6. **Keep tests focused** - one test should verify one behavior
7. **Avoid test interdependence** - tests should not rely on execution order
8. **Use table-driven tests** for testing multiple inputs/scenarios
9. **Mock external dependencies** - never call real APIs in unit tests
10. **Document complex test setups** with comments

## Common Patterns

### Testing Error Messages

```go
if out := ui.ErrorWriter.String(); !strings.Contains(out, "expected error") {
	t.Errorf("expected error message, got: %s", out)
}
```

### Testing Exit Codes

```go
if code := cmd.Run(args); code != 0 {
	t.Fatalf("expected exit 0, got %d", code)
}
```

### Testing Flag Parsing

```go
func TestCommandRequiresFlag(t *testing.T) {
	cmd := newCommand(cli.NewMockUi())

	if code := cmd.Run([]string{}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}
```

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [HashiCorp go-tfe Tests](https://github.com/hashicorp/go-tfe)
- [Testify Testing Library](https://github.com/stretchr/testify) (optional)
- [mitchellh/cli Testing](https://github.com/mitchellh/cli)
