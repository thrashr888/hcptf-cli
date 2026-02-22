#!/usr/bin/env bash
# scripts/api-coverage.sh - Compare API/go-tfe operation coverage against registered CLI commands.
#
# Usage: ./scripts/api-coverage.sh [summary|detail|missing]
#   summary - Show coverage table (default)
#   detail  - Show coverage table with missing operation list
#   missing - Show only resources with gaps

set -euo pipefail

MODE="${1:-summary}"
case "${MODE}" in
    summary|detail|missing) ;;
    *)
        echo "Usage: ./scripts/api-coverage.sh [summary|detail|missing]" >&2
        exit 1
        ;;
esac

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="${SCRIPT_DIR}/.."
COMMANDS_FILE="${REPO_ROOT}/command/commands.go"
CMD_DIR="${REPO_ROOT}/command"

if [[ ! -f "${COMMANDS_FILE}" ]]; then
    echo "Error: commands file not found at ${COMMANDS_FILE}" >&2
    exit 1
fi

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

if [[ ! -t 1 ]]; then
    RED='' GREEN='' YELLOW='' CYAN='' BOLD='' DIM='' NC=''
fi

declare -A CMD_SET=()
while IFS= read -r key; do
    [[ -z "${key}" ]] && continue
    CMD_SET["${key}"]=1
done < <(rg 'func\(\) \(cli.Command, error\)' "${COMMANDS_FILE}" | awk -F'"' '{print $2}')

# API resource definitions
# Format: "api_doc_name|command_prefix|operations"
# operations:
#   L=list C=create R=read U=update D=delete
#   plus custom ops (e.g. lock, list-org, add-workspace, outputs, query, override)
RESOURCES=(
    # Account & auth
    "account|account|C,R,U"

    # Organizations
    "organizations|organization|L,C,R,U,D"
    "organization-memberships|organizationmembership|L,C,R,D"
    "organization-tags|organizationtag|L,C,D"
    "organization-tokens|organizationtoken|L,C,R,D"

    # Workspaces
    "workspaces|workspace|L,C,R,U,D,lock,unlock,force-unlock"
    "workspace-variables|variable|L,C,U,D"
    "workspace-resources|workspaceresource|L,R"
    "workspace-tags|workspacetag|L,C,D"

    # Runs / plans / applies
    "runs|run|L,C,R,list-org,apply,discard,cancel,force-execute"
    "applies|apply|R"
    "plans|plan|R"
    "plan-exports|planexport|C,R,D"
    "cost-estimates|costestimate|R"

    # State
    "state-versions|state|L,R"
    "state-version-outputs|state|outputs"

    # Config versions / variable sets
    "configuration-versions|configversion|L,C,R"
    "variable-sets|variableset|L,C,R,U,D,apply-workspaces,remove-workspaces,apply-projects,remove-projects,apply-stacks,remove-stacks,list-workspace,list-project,update-workspaces,update-stacks"

    # Teams / access
    "teams|team|L,C,R,D"
    "team-access|teamaccess|L,C,R,U,D"
    "team-members|team|R"
    "team-tokens|teamtoken|L,C,R,D"

    # Projects
    "projects|project|L,C,R,U,D"
    "project-team-access|projectteamaccess|L,C,R,U,D"

    # Policies
    "policies|policy|L,C,R,U,D,upload,download"
    "policy-sets|policyset|L,C,R,U,D,add-policy,remove-policy,add-workspace,remove-workspace,add-workspace-exclusion,remove-workspace-exclusion,add-project,remove-project"
    "policy-checks|policycheck|L,R,override"
    "policy-evaluations|policyevaluation|L"
    "policy-set-params|policysetparameter|L,C,U,D"

    # Agents / pools
    "agents|agent|L,R"
    "agent-tokens|agentpool_token|L,C,D"

    # SSH / OAuth / notifications
    "ssh-keys|sshkey|L,C,R,U,D"
    "oauth-clients|oauthclient|L,C,R,U,D"
    "oauth-tokens|oauthtoken|L,R,U,D"
    "notification-configurations|notification|L,C,R,U,D"

    # Run tasks / triggers
    "run-tasks/run-tasks|runtask|L,C,R,U,D"
    "run-tasks/run-task-stages-and-results|runtask|R"
    "run-triggers|runtrigger|L,C,R,D"

    # Comments / audit
    "comments|comment|L,C,R"
    "audit-trails|audittrail|L,R"
    "audit-trails-tokens|audittrailtoken|L,C,R,D"

    # Registry
    "private-registry/modules|registrymodule|L,C,R,D"
    "private-registry/provider-versions-platforms|registryproviderplatform|C,R,D"
    "private-registry/providers|registryprovider|L,C,R,D"
    "private-registry/manage-provider-versions|registryproviderversion|C,R,D"
    "private-registry/gpg-keys|gpgkey|L,C,R,U,D"

    # VCS / health
    "vcs-events|vcsevent|L,R"
    "github-app-installations|githubapp|L,R"
    "assessment-results|assessmentresult|L,R"
    "change-requests|changerequest|L,C,R,U"

    # Stacks
    "stacks/stacks|stack|L,C,R,U,D"
    "stacks/stack-configurations|stackconfiguration|L,C,R,U,D"
    "stacks/stack-deployments|stackdeployment|L,C,R"
    "stacks/stack-states|stackstate|L,R"

    # OIDC / HYOK
    "hold-your-own-key/aws|awsoidc|C,R,U,D"
    "hold-your-own-key/azure|azureoidc|C,R,U,D"
    "hold-your-own-key/gcp|gcpoidc|C,R,U,D"
    "hold-your-own-key/vault-transit|vaultoidc|C,R,U,D"
    "hold-your-own-key/byok|hyok|L,C,R,U,D"
    "hold-your-own-key/key-management|hyokkey|C,R,D"

    # Queries / explorer
    "queries/run-query|queryrun|L"
    "queries/workspace-query|queryworkspace|L"
    "explorer|explorer|query"

    # User tokens / users
    "user-tokens|usertoken|L,C,R,D"
    "users|user|R"

    # Billing / metadata
    "subscriptions|subscription|L,R"
    "feature-sets|featureset|L"
    "ip-ranges|iprange|L"
    "no-code-provisioning|nocode|L,C,R,U"
    "stability-policy|stabilitypolicy|R"

    # Reserved tags
    "reserved-tag-keys|reservedtagkey|L,C,U,D"
)

op_label() {
    case "$1" in
        L) echo "list" ;;
        C) echo "create" ;;
        R) echo "read" ;;
        U) echo "update" ;;
        D) echo "delete" ;;
        *) echo "$1" ;;
    esac
}

resolve_prefix() {
    case "$1" in
        workspacetag) echo "workspace tag" ;;
        workspaceresource) echo "workspace resource" ;;
        stackconfiguration) echo "stack configuration" ;;
        stackdeployment) echo "stack deployment" ;;
        stackstate) echo "stack state" ;;
        registrymodule) echo "registry module" ;;
        registryprovider) echo "registry provider" ;;
        registryproviderversion) echo "registry provider version" ;;
        registryproviderplatform) echo "registry provider platform" ;;
        agentpool_token) echo "agentpool token" ;;
        *) echo "$1" ;;
    esac
}

has_registered_command() {
    local cmd="$1"
    [[ -n "${CMD_SET[${cmd}]+x}" ]]
}

has_command() {
    local raw_prefix="$1"
    local action="$2"
    local prefix
    prefix="$(resolve_prefix "${raw_prefix}")"

    local candidates=()

    # Main CRUD mapping.
    case "${prefix}:${action}" in
        agentpool\ token:list) candidates+=("agentpool token-list") ;;
        agentpool\ token:create) candidates+=("agentpool token-create") ;;
        agentpool\ token:delete) candidates+=("agentpool token-delete") ;;
        *) candidates+=("${prefix} ${action}") ;;
    esac

    # Read aliases.
    if [[ "${action}" == "read" ]]; then
        candidates+=("${prefix} show")
    fi

    # Custom operation aliases.
    case "${prefix}:${action}" in
        workspace\ tag:create) candidates+=("workspace tag add") ;;
        workspace\ tag:delete) candidates+=("workspace tag remove") ;;
        workspace:force-unlock) candidates+=("workspace force-unlock") ;;
        variableset:apply-workspaces) candidates+=("variableset apply") ;;
        variableset:remove-workspaces) candidates+=("variableset remove") ;;
        variableset:apply-projects) candidates+=("variableset apply") ;;
        variableset:remove-projects) candidates+=("variableset remove") ;;
        variableset:apply-stacks) candidates+=("variableset apply") ;;
        variableset:remove-stacks) candidates+=("variableset remove") ;;
        explorer:query) candidates+=("explorer query") ;;
        policycheck:override) candidates+=("policycheck override") ;;
        state:outputs) candidates+=("state outputs") ;;
    esac

    local candidate
    for candidate in "${candidates[@]}"; do
        if has_registered_command "${candidate}"; then
            return 0
        fi
    done

    # Fallback for command files that may exist without registration.
    local file_prefix="${raw_prefix// /}"
    if [[ -f "${CMD_DIR}/${file_prefix}_${action}.go" ]]; then
        return 0
    fi

    return 1
}

total_resources=0
covered_resources=0
fully_covered_resources=0
total_ops=0
covered_ops=0
missing_resources=()
missing_ops=()

W_RES=42
W_OPS=28
W_MISS=28
W_COV=8

print_header() {
    printf "${BOLD}%-${W_RES}s │ %-${W_OPS}s │ %-${W_MISS}s │ %s${NC}\n" \
        "API Resource" "Covered Ops" "Missing Ops" "Coverage"
    printf "%-${W_RES}s─┼─%-${W_OPS}s─┼─%-${W_MISS}s─┼─%s\n" \
        "$(printf '─%.0s' $(seq 1 ${W_RES}))" \
        "$(printf '─%.0s' $(seq 1 ${W_OPS}))" \
        "$(printf '─%.0s' $(seq 1 ${W_MISS}))" \
        "$(printf '─%.0s' $(seq 1 ${W_COV}))"
}

print_row() {
    local resource="$1" found_str="$2" missing_str="$3" pct="$4"

    local color="${RED}"
    if [[ "${pct}" -eq 100 ]]; then
        color="${GREEN}"
    elif [[ "${pct}" -ge 50 ]]; then
        color="${YELLOW}"
    fi

    printf "%-${W_RES}s │ ${GREEN}%-${W_OPS}s${NC} │ ${RED}%-${W_MISS}s${NC} │ ${color}%3d%%${NC}\n" \
        "${resource}" "${found_str}" "${missing_str}" "${pct}"
}

echo ""
echo -e "${BOLD}TFC API Coverage Report${NC}"
echo -e "${DIM}Derived from command registry + resource operation expectations${NC}"
echo ""

print_header

for entry in "${RESOURCES[@]}"; do
    IFS='|' read -r api_doc cli_prefix ops_csv <<< "${entry}"

    IFS=',' read -ra ops <<< "${ops_csv}"
    resource_total=${#ops[@]}
    resource_found=0
    found_labels=()
    missing_labels=()

    for op in "${ops[@]}"; do
        action="$(op_label "${op}")"
        total_ops=$((total_ops + 1))

        if has_command "${cli_prefix}" "${action}"; then
            resource_found=$((resource_found + 1))
            covered_ops=$((covered_ops + 1))
            found_labels+=("${action}")
        else
            missing_labels+=("${action}")
            missing_ops+=("${api_doc}:${action}")
        fi
    done

    total_resources=$((total_resources + 1))
    if [[ "${resource_found}" -gt 0 ]]; then
        covered_resources=$((covered_resources + 1))
    else
        missing_resources+=("${api_doc}")
    fi
    if [[ "${resource_found}" -eq "${resource_total}" ]]; then
        fully_covered_resources=$((fully_covered_resources + 1))
    fi

    if [[ "${resource_total}" -gt 0 ]]; then
        pct=$((resource_found * 100 / resource_total))
    else
        pct=0
    fi

    found_str="$(IFS=' '; echo "${found_labels[*]:-}")"
    miss_str="$(IFS=' '; echo "${missing_labels[*]:-}")"

    case "${MODE}" in
        missing)
            if [[ "${pct}" -lt 100 ]]; then
                print_row "${api_doc}" "${found_str}" "${miss_str}" "${pct}"
            fi
            ;;
        *)
            print_row "${api_doc}" "${found_str}" "${miss_str}" "${pct}"
            ;;
    esac
done

echo ""
echo -e "${BOLD}Summary${NC}"

if [[ "${total_resources}" -gt 0 ]]; then
    res_pct=$((covered_resources * 100 / total_resources))
else
    res_pct=0
fi
if [[ "${total_ops}" -gt 0 ]]; then
    ops_pct=$((covered_ops * 100 / total_ops))
else
    ops_pct=0
fi

echo -e "  Resources with any CLI coverage: ${BOLD}${covered_resources}/${total_resources}${NC} (${res_pct}%)"
echo -e "  Resources fully covered:         ${BOLD}${fully_covered_resources}/${total_resources}${NC}"
echo -e "  Total operations covered:        ${BOLD}${covered_ops}/${total_ops}${NC} (${ops_pct}%)"

if [[ ${#missing_resources[@]} -gt 0 ]]; then
    echo ""
    echo -e "${BOLD}Resources with no covered operations:${NC}"
    for r in "${missing_resources[@]}"; do
        echo -e "  ${RED}✗${NC} ${r}"
    done
fi

if [[ "${MODE}" == "detail" && ${#missing_ops[@]} -gt 0 ]]; then
    echo ""
    echo -e "${BOLD}All missing operations:${NC}"
    for op in "${missing_ops[@]}"; do
        echo -e "  ${RED}✗${NC} ${op}"
    done
fi

echo ""
