# Test Coverage Report

## Current Coverage Summary

| Package | Coverage | Status |
|---------|----------|--------|
| internal/client | 92.9% | ✅ Excellent |
| internal/config | 74.0% | ✅ Good |
| internal/router | 73.2% | ✅ Good |
| internal/output | 60.9% | ⚠️  Fair |
| command | 12.8% | ❌ Needs Improvement |

## Package Details

### internal/client (92.9%)
**Status**: ✅ Well tested

**Test file**: `internal/client/client_test.go`

**Coverage includes**:
- Client creation with various token sources
- Environment variable precedence
- Invalid address handling
- Context creation
- Address retrieval

**Missing**:
- Custom HTTP client configuration (minor)

### internal/config (74.0%)
**Status**: ✅ Good coverage

**Test file**: `internal/config/config_test.go`

**Coverage includes**:
- Config file parsing
- Credential management
- Token precedence (env vars vs config)
- HCL parsing
- Default values

**Missing**:
- Edge cases in config file parsing
- Some error paths

### internal/router (73.2%)
**Status**: ✅ Good coverage

**Test file**: `internal/router/router_test.go`

**Coverage includes**:
- URL-style argument translation
- Known command detection
- All major routing patterns (org, workspace, runs)
- Edge cases

**Missing**:
- Context-based validation (optional feature)

### internal/output (60.9%)
**Status**: ⚠️ Fair coverage

**Test file**: `internal/output/output_test.go`

**Coverage includes**:
- Table formatting
- JSON formatting
- Row/column handling

**Missing**:
- Error handling edge cases
- Complex table scenarios

### command (12.8%)
**Status**: ❌ Needs significant improvement

**Test coverage**: 42 out of 254 command files have tests (16.5%)

**Commands WITH tests** (42 files):
- Workspace commands (6): list, create, read, update, delete, partial coverage
- Run commands (6): list, create, show, apply, discard, cancel
- Variable commands (3): create, update, delete
- Run Trigger commands (4): list, create, read, delete
- Plan commands (2): read, logs
- Apply commands (2): read, logs
- Config Version commands (2): list, read
- GPG Key commands (5): create, list, read, update, delete
- Registry Provider commands (4): list, create, read, delete
- Workspace Resource commands (2): list, read
- VCS Event commands (2): list, read
- Plan Export commands (3): create, download, read

**Commands WITHOUT tests** (212 files):
Priority for testing:

**High Priority** (Core Operations):
- Organization commands (5): list, create, read, update, delete
- State commands (3): list, read, outputs
- Project commands (5): list, create, read, update, delete
- Team commands (7): list, create, read, update, delete, add-member, remove-member
- Policy commands (6): list, create, read, update, delete, upload
- Policy Set commands (6): list, create, read, update, delete, add-policy

**Medium Priority** (Important Features):
- Variable Set commands (10): all operations
- Agent Pool commands (8): all operations
- Run Task commands (7): all operations
- OAuth Client commands (5): all operations
- Team Access commands (5): all operations
- Notification commands (5): all operations
- SSH Key commands (6): all operations

**Lower Priority** (Enterprise/Advanced):
- Stack commands (all)
- Audit Trail commands (all)
- OIDC commands (all 4 providers)
- HYOK commands (all)
- Assessment commands (all)
- Change Request commands (all)
- Comment commands (all)
- Organization Membership commands (all)
- Organization Tag commands (all)
- Policy Check/Evaluation commands (all)

## Testing Progress Tracking

### Phase 1: Foundation ✅ COMPLETE
- [x] internal/client tests
- [x] internal/config tests
- [x] internal/router tests
- [x] internal/output tests
- [x] Test helpers and mocks
- [x] Testing documentation

### Phase 2: Core Commands (Target: 40% coverage)
- [ ] Organization commands (5 files)
- [ ] State commands (3 files)
- [ ] Project commands (5 files)
- [ ] Team commands (7 files)

### Phase 3: Important Features (Target: 60% coverage)
- [ ] Variable Set commands (10 files)
- [ ] Agent Pool commands (8 files)
- [ ] Run Task commands (7 files)
- [ ] OAuth Client commands (5 files)
- [ ] Team Access commands (5 files)

### Phase 4: Advanced Features (Target: 75% coverage)
- [ ] Policy commands (6 files)
- [ ] Policy Set commands (6 files)
- [ ] Notification commands (5 files)
- [ ] SSH Key commands (6 files)

### Phase 5: Enterprise Features (Target: 80%+ coverage)
- [ ] Stack commands
- [ ] Audit Trail commands
- [ ] OIDC commands
- [ ] HYOK commands
- [ ] Advanced Policy commands

## Test Quality Metrics

### Good Test Patterns
✅ Using dependency injection with mock services
✅ Table-driven tests for multiple scenarios
✅ Testing both success and error paths
✅ Testing output formats (JSON and table)
✅ Testing flag validation
✅ Proper use of test helpers

### Areas for Improvement
⚠️ Need more integration tests
⚠️ Command coverage still low (12.8%)
⚠️ Missing edge case coverage
⚠️ Need performance/benchmark tests
⚠️ Need end-to-end tests

## Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/client/...
go test ./command/...

# Run with race detector
go test -race ./...
```

## Goals

| Timeframe | Target | Status |
|-----------|--------|--------|
| Current | Internal packages > 70% | ✅ Achieved (79% avg) |
| Short-term | Command package > 40% | 🔄 In Progress (12.8%) |
| Mid-term | Command package > 60% | ⏳ Planned |
| Long-term | Overall > 75% | ⏳ Planned |

## Next Steps

1. Add tests for Organization commands (high impact, foundational)
2. Add tests for State commands (frequently used)
3. Add tests for Project commands (organizational structure)
4. Add tests for Team commands (permission management)
5. Create integration test framework
6. Add benchmark tests for performance-critical paths

## Contributing

When adding new commands:
1. Always include corresponding test file
2. Cover all required flags
3. Test error handling
4. Test output formats (table and JSON)
5. Use existing test patterns from `test_helpers_test.go`

See [TESTING.md](TESTING.md) for detailed testing guidelines.
