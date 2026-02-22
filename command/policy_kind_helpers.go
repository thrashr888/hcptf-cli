package command

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

func parsePolicyKind(value string) (tfe.PolicyKind, error) {
	switch value {
	case "", "sentinel":
		return tfe.Sentinel, nil
	case "opa":
		return tfe.OPA, nil
	default:
		return "", fmt.Errorf("invalid policy kind %q: must be 'sentinel' or 'opa'", value)
	}
}
