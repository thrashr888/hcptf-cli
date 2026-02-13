#!/bin/bash
# Coverage analysis script for hcptf-cli
# Usage: ./scripts/coverage.sh [html|report|summary]

set -e

MODE="${1:-summary}"

# Colors
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "Running tests with coverage..."
go test -coverprofile=coverage.out ./... > /dev/null 2>&1

case "$MODE" in
  html)
    echo "Generating HTML coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓${NC} Coverage report generated: coverage.html"
    echo "Open in browser with: open coverage.html"
    ;;

  report)
    echo ""
    echo "========================================"
    echo "        Coverage Report by Package"
    echo "========================================"
    echo ""

    go test -cover ./... 2>&1 | grep -E 'ok|FAIL' | while read line; do
      if [[ $line =~ coverage:\ ([0-9]+\.[0-9]+)% ]]; then
        coverage="${BASH_REMATCH[1]}"
        package=$(echo "$line" | awk '{print $2}')

        # Color code based on coverage
        if (( $(echo "$coverage >= 80" | bc -l) )); then
          color=$GREEN
          status="✓ Excellent"
        elif (( $(echo "$coverage >= 60" | bc -l) )); then
          color=$BLUE
          status="○ Good"
        elif (( $(echo "$coverage >= 40" | bc -l) )); then
          color=$YELLOW
          status="△ Fair"
        else
          color=$RED
          status="✗ Needs Work"
        fi

        printf "${color}%-60s %6.1f%%  %s${NC}\n" "$package" "$coverage" "$status"
      fi
    done

    echo ""
    echo "========================================"
    echo "        Untested Functions"
    echo "========================================"
    echo ""

    go tool cover -func=coverage.out | grep "0.0%" | head -20

    count=$(go tool cover -func=coverage.out | grep -c "0.0%" || true)
    echo ""
    echo "Total untested functions: $count"
    ;;

  summary)
    echo ""
    echo "========================================"
    echo "        Coverage Summary"
    echo "========================================"
    echo ""

    # Overall coverage
    total=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
    echo -e "Overall Coverage:      ${BLUE}$total%${NC}"
    echo ""

    # Package breakdown
    echo "Package Breakdown:"
    echo ""

    # Internal packages
    go test -cover ./... 2>&1 | grep -E 'internal' | while read line; do
      if [[ $line =~ coverage:\ ([0-9]+\.[0-9]+)% ]]; then
        coverage="${BASH_REMATCH[1]}"
        package=$(echo "$line" | awk '{print $2}' | sed 's/.*\///')

        if (( $(echo "$coverage >= 80" | bc -l) )); then
          color=$GREEN
        elif (( $(echo "$coverage >= 60" | bc -l) )); then
          color=$BLUE
        elif (( $(echo "$coverage >= 40" | bc -l) )); then
          color=$YELLOW
        else
          color=$RED
        fi

        printf "  ${color}%-20s %6.1f%%${NC}\n" "$package" "$coverage"
      fi
    done

    # Command package
    cmd_coverage=$(go test -cover ./command/... 2>&1 | grep coverage | awk '{print $5}' | sed 's/%//')
    if (( $(echo "$cmd_coverage >= 80" | bc -l) )); then
      color=$GREEN
    elif (( $(echo "$cmd_coverage >= 60" | bc -l) )); then
      color=$BLUE
    elif (( $(echo "$cmd_coverage >= 40" | bc -l) )); then
      color=$YELLOW
    else
      color=$RED
    fi
    printf "  ${color}%-20s %6.1f%%${NC}\n" "command" "$cmd_coverage"

    echo ""
    echo "========================================"
    echo "        Coverage Goals"
    echo "========================================"
    echo ""

    # Calculate progress to goals
    goal_50=50
    goal_75=75
    goal_80=80

    if (( $(echo "$total >= $goal_80" | bc -l) )); then
      echo -e "  ${GREEN}✓${NC} Exceeded 80% goal! 🎉"
    elif (( $(echo "$total >= $goal_75" | bc -l) )); then
      echo -e "  ${GREEN}✓${NC} Exceeded 75% goal!"
      remaining=$(echo "$goal_80 - $total" | bc)
      echo -e "  ${YELLOW}○${NC} Next goal: 80% (need ${remaining}% more)"
    elif (( $(echo "$total >= $goal_50" | bc -l) )); then
      echo -e "  ${GREEN}✓${NC} Exceeded 60% goal!"
      remaining=$(echo "$goal_75 - $total" | bc)
      echo -e "  ${YELLOW}○${NC} Next goal: 75% (need ${remaining}% more)"
    else
      remaining=$(echo "$goal_50 - $total" | bc)
      echo -e "  ${YELLOW}○${NC} Target: 60% (need ${remaining}% more)"
    fi

    echo ""
    echo "Run './scripts/coverage.sh report' for detailed breakdown"
    echo "Run './scripts/coverage.sh html' to view coverage in browser"
    ;;

  *)
    echo "Usage: $0 [html|report|summary]"
    echo ""
    echo "  summary  - Show coverage summary (default)"
    echo "  report   - Show detailed coverage by package and untested functions"
    echo "  html     - Generate HTML coverage report"
    exit 1
    ;;
esac

echo ""
