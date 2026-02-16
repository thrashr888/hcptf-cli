package command

import "testing"

func TestCommandsNamespaceNormalizationCanonicalOnly(t *testing.T) {
	commandMap := Commands(&Meta{})

	resolveCommand := func(name string) {
		t.Helper()

		factory, ok := commandMap[name]
		if !ok {
			t.Fatalf("missing command key %q", name)
		}

		_, err := factory()
		if err != nil {
			t.Fatalf("failed to construct command %q: %v", name, err)
		}
	}

	canonicalKeys := []string{
		"audittrail token",
		"audittrail token list",
		"audittrail token create",
		"audittrail token read",
		"audittrail token delete",
		"organization member",
		"organization member read",
		"organization membership",
		"organization membership list",
		"organization membership create",
		"organization membership read",
		"organization membership delete",
		"organization token",
		"organization token list",
		"organization token create",
		"organization token read",
		"organization token delete",
		"organization tag",
		"organization tag list",
		"organization tag create",
		"organization tag delete",
		"policyset parameter",
		"policyset parameter list",
		"policyset parameter create",
		"policyset parameter update",
		"policyset parameter delete",
		"policyset outcome",
		"policyset outcome list",
		"policyset outcome read",
		"project teamaccess",
		"project teamaccess list",
		"project teamaccess create",
		"project teamaccess read",
		"project teamaccess update",
		"project teamaccess delete",
		"team access",
		"team access list",
		"team access create",
		"team access read",
		"team access update",
		"team access delete",
		"team token",
		"team token list",
		"team token create",
		"team token read",
		"team token delete",
		"user token",
		"user token list",
		"user token create",
		"user token read",
		"user token delete",
		"workspace resource",
		"workspace resource list",
		"workspace resource read",
		"workspace tag",
		"workspace tag list",
		"workspace tag add",
		"workspace tag remove",
	}

	for _, key := range canonicalKeys {
		resolveCommand(key)
	}

	legacyKeys := []string{
		"audittrailtoken",
		"audittrailtoken list",
		"audittrailtoken create",
		"audittrailtoken read",
		"audittrailtoken delete",
		"organizationmember",
		"organizationmember read",
		"organizationmembership",
		"organizationmembership list",
		"organizationmembership create",
		"organizationmembership read",
		"organizationmembership delete",
		"organizationtoken",
		"organizationtoken list",
		"organizationtoken create",
		"organizationtoken read",
		"organizationtoken delete",
		"organizationtag",
		"organizationtag list",
		"organizationtag create",
		"organizationtag delete",
		"policysetparameter",
		"policysetparameter list",
		"policysetparameter create",
		"policysetparameter update",
		"policysetparameter delete",
		"policysetoutcome",
		"policysetoutcome list",
		"policysetoutcome read",
		"projectteamaccess",
		"projectteamaccess list",
		"projectteamaccess create",
		"projectteamaccess read",
		"projectteamaccess update",
		"projectteamaccess delete",
		"teamaccess",
		"teamaccess list",
		"teamaccess create",
		"teamaccess read",
		"teamaccess update",
		"teamaccess delete",
		"teamtoken",
		"teamtoken list",
		"teamtoken create",
		"teamtoken read",
		"teamtoken delete",
		"workspaceresource",
		"workspaceresource list",
		"workspaceresource read",
		"workspacetag",
		"workspacetag list",
		"workspacetag add",
		"workspacetag remove",
		"usertoken",
		"usertoken list",
		"usertoken create",
		"usertoken read",
		"usertoken delete",
	}

	for _, key := range legacyKeys {
		if _, ok := commandMap[key]; ok {
			t.Fatalf("legacy command key should not exist: %q", key)
		}
	}
}
