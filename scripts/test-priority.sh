#!/bin/bash
# Script to identify high-priority commands for testing
# Analyzes which commands are most important to test based on usage patterns

set -e

echo "========================================"
echo "   Test Priority Analysis"
echo "========================================"
echo ""

# High priority command groups (frequently used, core functionality)
HIGH_PRIORITY=(
  "organization"
  "workspace"
  "run"
  "state"
  "variable"
  "project"
  "team"
)

# Medium priority
MEDIUM_PRIORITY=(
  "variableset"
  "agentpool"
  "runtask"
  "oauthclient"
  "policy"
  "policyset"
  "sshkey"
  "notification"
)

echo "HIGH PRIORITY (Core Operations)"
echo "--------------------------------"
for cmd in "${HIGH_PRIORITY[@]}"; do
  # Count files
  files=$(ls command/${cmd}_*.go 2>/dev/null | wc -l | tr -d ' ')
  tests=$(ls command/${cmd}_*_test.go 2>/dev/null | wc -l | tr -d ' ')

  if [ "$tests" -gt 0 ]; then
    status="✓ Has tests ($tests/$files)"
    color="\033[0;32m"
  else
    status="✗ No tests (0/$files)"
    color="\033[0;31m"
  fi

  printf "${color}%-20s %s\033[0m\n" "$cmd" "$status"
done

echo ""
echo "MEDIUM PRIORITY (Important Features)"
echo "-------------------------------------"
for cmd in "${MEDIUM_PRIORITY[@]}"; do
  files=$(ls command/${cmd}_*.go 2>/dev/null | wc -l | tr -d ' ')
  tests=$(ls command/${cmd}_*_test.go 2>/dev/null | wc -l | tr -d ' ')

  if [ "$tests" -gt 0 ]; then
    status="✓ Has tests ($tests/$files)"
    color="\033[0;32m"
  elif [ "$files" -eq 0 ]; then
    continue
  else
    status="✗ No tests (0/$files)"
    color="\033[0;31m"
  fi

  printf "${color}%-20s %s\033[0m\n" "$cmd" "$status"
done

echo ""
echo "========================================"
echo "   Recommended Testing Order"
echo "========================================"
echo ""
echo "Based on priority and current coverage:"
echo ""

# Find high priority commands without tests
echo "1. HIGH PRIORITY - Add tests first:"
for cmd in "${HIGH_PRIORITY[@]}"; do
  tests=$(ls command/${cmd}_*_test.go 2>/dev/null | wc -l | tr -d ' ')
  files=$(ls command/${cmd}_*.go 2>/dev/null | wc -l | tr -d ' ')

  if [ "$tests" -eq 0 ] && [ "$files" -gt 0 ]; then
    printf "   - %-20s (%d files to test)\n" "$cmd" "$files"
  fi
done

echo ""
echo "2. MEDIUM PRIORITY - Add tests next:"
for cmd in "${MEDIUM_PRIORITY[@]}"; do
  tests=$(ls command/${cmd}_*_test.go 2>/dev/null | wc -l | tr -d ' ')
  files=$(ls command/${cmd}_*.go 2>/dev/null | wc -l | tr -d ' ')

  if [ "$tests" -eq 0 ] && [ "$files" -gt 0 ]; then
    printf "   - %-20s (%d files to test)\n" "$cmd" "$files"
  fi
done

echo ""
echo "========================================"
echo "   Quick Stats"
echo "========================================"
echo ""

total_commands=$(ls command/*_*.go 2>/dev/null | grep -v _test.go | wc -l | tr -d ' ')
total_tests=$(ls command/*_test.go 2>/dev/null | wc -l | tr -d ' ')
coverage_pct=$(echo "scale=1; ($total_tests / $total_commands) * 100" | bc)

echo "Total command files: $total_commands"
echo "Total test files:    $total_tests"
echo "Test coverage:       ${coverage_pct}%"
echo ""
