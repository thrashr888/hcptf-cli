#!/usr/bin/env bash
# scripts/api-coverage.sh - Compare TFC API docs against CLI command coverage
#
# Usage: ./scripts/api-coverage.sh [summary|detail|missing]
#   summary - Show coverage table (default)
#   detail  - Show coverage table with per-resource operation details
#   missing - Show only missing resources and operations

set -euo pipefail

MODE="${1:-summary}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CMD_DIR="${SCRIPT_DIR}/../command"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# Disable colors if not a terminal
if [ ! -t 1 ]; then
    RED='' GREEN='' YELLOW='' CYAN='' BOLD='' DIM='' NC=''
fi

# ─────────────────────────────────────────────────────────────────────────────
# API resource definitions
# Format: "api_doc_name|cli_prefix|operations"
#
# api_doc_name: filename from web-unified-docs/.../api-docs/ (without .mdx)
# cli_prefix:   prefix used in command/<prefix>_<action>.go
# operations:   comma-separated list of expected CRUD operations
#               L=list, C=create, R=read, U=update, D=delete, plus custom ops
#
# Source: https://github.com/hashicorp/web-unified-docs/tree/main/content/terraform-docs-common/docs/cloud-docs/api-docs
# ─────────────────────────────────────────────────────────────────────────────

RESOURCES=(
    # Account & Authentication
    "account|account|C,R,U"

    # Organizations
    "organizations|organization|L,C,R,U,D"
    "organization-memberships|organizationmembership|L,C,R,D"
    "organization-tags|organizationtag|L,D"
    "organization-tokens|organizationtoken|L,C,R,D"

    # Workspaces
    "workspaces|workspace|L,C,R,U,D"
    "workspace-variables|variable|L,C,U,D"
    "workspace-resources|workspaceresource|L,R"

    # Runs
    "runs|run|L,C,R"
    "applies|apply|R"
    "plans|plan|R"
    "plan-exports|planexport|C,R"
    "cost-estimates|costestimate|R"

    # State
    "state-versions|state|L,R"
    "state-version-outputs|state|R"

    # Configuration
    "configuration-versions|configversion|L,C,R"

    # Variables
    "variable-sets|variableset|L,C,R,U,D"

    # Teams & Access
    "teams|team|L,C,R,D"
    "team-access|teamaccess|L,C,R,U,D"
    "team-members|team|R"
    "team-tokens|teamtoken|L,C,R,D"

    # Projects
    "projects|project|L,C,R,U,D"
    "project-team-access|projectteamaccess|L,C,R,U,D"

    # Policies
    "policies|policy|L,C,R,U,D"
    "policy-sets|policyset|L,C,R,U,D"
    "policy-checks|policycheck|L,R"
    "policy-evaluations|policyevaluation|L"
    "policy-set-params|policysetparameter|L,C,U,D"

    # Agents
    "agents|agent|L,R"
    "agent-tokens|agentpool_token|L,C,D"

    # SSH Keys
    "ssh-keys|sshkey|L,C,R,U,D"

    # OAuth
    "oauth-clients|oauthclient|L,C,R,U,D"
    "oauth-tokens|oauthtoken|L,R,U"

    # Notifications
    "notification-configurations|notification|L,C,R,U,D"

    # Run Tasks
    "run-tasks/run-tasks|runtask|L,C,R,U,D"
    "run-tasks/run-task-stages-and-results|runtask|R"

    # Run Triggers
    "run-triggers|runtrigger|L,C,R,D"

    # Comments
    "comments|comment|L,C,R"

    # Audit Trails
    "audit-trails|audittrail|L,R"
    "audit-trails-tokens|audittrailtoken|L,C,R,D"

    # Registry - Modules
    "private-registry/modules|registrymodule|L,C,R,D"
    "private-registry/provider-versions-platforms|registryproviderversion|C,R,D"
    "private-registry/providers|registryprovider|L,C,R,D"
    "private-registry/gpg-keys|gpgkey|L,C,R,U,D"

    # VCS
    "vcs-events|vcsevent|L,R"
    "github-app-installations|githubapp|L,R"

    # Assessment / Health
    "assessment-results|assessmentresult|L,R"

    # Change Requests
    "change-requests|changerequest|L,C,R,U"

    # Stacks
    "stacks/stacks|stack|L,C,R,U,D"
    "stacks/stack-configurations|stackconfiguration|L,C,R,U,D"
    "stacks/stack-deployments|stackdeployment|L,C,R"
    "stacks/stack-states|stackstate|L,R"

    # OIDC configurations
    "hold-your-own-key/aws|awsoidc|C,R,U,D"
    "hold-your-own-key/azure|azureoidc|C,R,U,D"
    "hold-your-own-key/gcp|gcpoidc|C,R,U,D"
    "hold-your-own-key/vault-transit|vaultoidc|C,R,U,D"
    "hold-your-own-key/byok|hyok|L,C,R,U,D"
    "hold-your-own-key/key-management|hyokkey|C,R,D"

    # Explorer / Queries
    "queries/run-query|queryrun|L"
    "queries/workspace-query|queryworkspace|L"

    # User Tokens
    "user-tokens|usertoken|L,C,R,D"
    "users|account|R,U"

    # Subscriptions & Billing (read-only / not typically CLI)
    "subscriptions|subscription|R"
    "invoices|invoice|L,R"
    "feature-sets|featureset|L"

    # IP Ranges (special endpoint)
    "ip-ranges|iprange|R"

    # Reserved Tag Keys
    "reserved-tag-keys|reservedtagkey|L,C,D"

    # No-Code Provisioning
    "no-code-provisioning|nocode|L,C,R,U,D"

    # Stability Policy
    "stability-policy|stabilitypolicy|R"
)

# Operation label map
op_label() {
    case "$1" in
        L) echo "list";;
        C) echo "create";;
        R) echo "read";;
        U) echo "update";;
        D) echo "delete";;
        *) echo "$1";;
    esac
}

# Check if a command file exists for a given prefix and action
has_command() {
    local prefix="$1" action="$2"
    # Check for exact file match
    if [ -f "${CMD_DIR}/${prefix}_${action}.go" ]; then
        return 0
    fi
    # Special cases
    case "${prefix}_${action}" in
        account_read) [ -f "${CMD_DIR}/account_show.go" ] && return 0;;
        organization_read) [ -f "${CMD_DIR}/organization_show.go" ] && return 0;;
        team_read) [ -f "${CMD_DIR}/team_show.go" ] && return 0;;
        apply_read) [ -f "${CMD_DIR}/apply_read.go" ] && return 0;;
        run_read) [ -f "${CMD_DIR}/run_show.go" ] && return 0;;
    esac
    return 1
}

# ─────────────────────────────────────────────────────────────────────────────
# Main
# ─────────────────────────────────────────────────────────────────────────────

total_resources=0
covered_resources=0
total_ops=0
covered_ops=0
missing_resources=()
missing_ops=()

# Column widths
W_RES=38
W_OPS=22
W_MISS=22
W_COV=8

print_header() {
    printf "${BOLD}%-${W_RES}s │ %-${W_OPS}s │ %-${W_MISS}s │ %s${NC}\n" \
        "API Resource" "CLI Commands" "Missing" "Coverage"
    printf "%-${W_RES}s─┼─%-${W_OPS}s─┼─%-${W_MISS}s─┼─%s\n" \
        "$(printf '─%.0s' $(seq 1 $W_RES))" \
        "$(printf '─%.0s' $(seq 1 $W_OPS))" \
        "$(printf '─%.0s' $(seq 1 $W_MISS))" \
        "$(printf '─%.0s' $(seq 1 $W_COV))"
}

print_row() {
    local resource="$1" found_str="$2" missing_str="$3" pct="$4"

    local color="$RED"
    if [ "$pct" -eq 100 ]; then
        color="$GREEN"
    elif [ "$pct" -ge 50 ]; then
        color="$YELLOW"
    fi

    printf "%-${W_RES}s │ ${GREEN}%-${W_OPS}s${NC} │ ${RED}%-${W_MISS}s${NC} │ ${color}%3d%%${NC}\n" \
        "$resource" "$found_str" "$missing_str" "$pct"
}

echo ""
echo -e "${BOLD}TFC API Coverage Report${NC}"
echo -e "${DIM}Comparing Terraform Cloud API docs against CLI commands${NC}"
echo ""

print_header

for entry in "${RESOURCES[@]}"; do
    IFS='|' read -r api_doc cli_prefix ops_csv <<< "$entry"

    IFS=',' read -ra ops <<< "$ops_csv"
    resource_total=${#ops[@]}
    resource_found=0
    found_labels=()
    missing_labels=()

    for op in "${ops[@]}"; do
        action=$(op_label "$op")
        total_ops=$((total_ops + 1))

        if has_command "$cli_prefix" "$action"; then
            resource_found=$((resource_found + 1))
            covered_ops=$((covered_ops + 1))
            found_labels+=("$op")
        else
            missing_labels+=("$op")
            missing_ops+=("${api_doc}:${action}")
        fi
    done

    total_resources=$((total_resources + 1))
    if [ "$resource_found" -gt 0 ]; then
        covered_resources=$((covered_resources + 1))
    else
        missing_resources+=("$api_doc")
    fi

    if [ "$resource_total" -gt 0 ]; then
        pct=$((resource_found * 100 / resource_total))
    else
        pct=0
    fi

    found_str=$(IFS=' '; echo "${found_labels[*]:-}")
    miss_str=$(IFS=' '; echo "${missing_labels[*]:-}")

    case "$MODE" in
        missing)
            if [ "$pct" -lt 100 ]; then
                print_row "$api_doc" "$found_str" "$miss_str" "$pct"
            fi
            ;;
        *)
            print_row "$api_doc" "$found_str" "$miss_str" "$pct"
            ;;
    esac
done

# Summary
echo ""
echo -e "${BOLD}Summary${NC}"

if [ "$total_resources" -gt 0 ]; then
    res_pct=$((covered_resources * 100 / total_resources))
else
    res_pct=0
fi
if [ "$total_ops" -gt 0 ]; then
    ops_pct=$((covered_ops * 100 / total_ops))
else
    ops_pct=0
fi

echo -e "  Resources with any CLI coverage: ${BOLD}${covered_resources}/${total_resources}${NC} (${res_pct}%)"
echo -e "  Total operations covered:        ${BOLD}${covered_ops}/${total_ops}${NC} (${ops_pct}%)"

if [ ${#missing_resources[@]} -gt 0 ]; then
    echo ""
    echo -e "${BOLD}Resources with no CLI commands:${NC}"
    for r in "${missing_resources[@]}"; do
        echo -e "  ${RED}✗${NC} ${r}"
    done
fi

if [ "$MODE" = "detail" ] && [ ${#missing_ops[@]} -gt 0 ]; then
    echo ""
    echo -e "${BOLD}All missing operations:${NC}"
    for op in "${missing_ops[@]}"; do
        echo -e "  ${RED}✗${NC} ${op}"
    done
fi

echo ""
