package command

import (
	"github.com/mitchellh/cli"
	"strings"
)

// Commands returns the mapping of CLI commands
func Commands(meta *Meta) map[string]cli.CommandFactory {
	commands := map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				Meta: *meta,
			}, nil
		},

		// Context help commands (internal, for URL-like arg support)
		"organization:context": func() (cli.Command, error) {
			return &OrganizationContextCommand{
				Meta: *meta,
			}, nil
		},
		"workspace:context": func() (cli.Command, error) {
			return &WorkspaceContextCommand{
				Meta: *meta,
			}, nil
		},

		// Authentication commands
		"login": func() (cli.Command, error) {
			return &LoginCommand{
				Meta: *meta,
			}, nil
		},
		"whoami": func() (cli.Command, error) {
			return &WhoAmICommand{
				Meta: *meta,
			}, nil
		},
		"logout": func() (cli.Command, error) {
			return &LogoutCommand{
				Meta: *meta,
			}, nil
		},

		// Account commands
		"account create": func() (cli.Command, error) {
			return &AccountCreateCommand{
				Meta: *meta,
			}, nil
		},
		"account show": func() (cli.Command, error) {
			return &AccountShowCommand{
				Meta: *meta,
			}, nil
		},
		"account update": func() (cli.Command, error) {
			return &AccountUpdateCommand{
				Meta: *meta,
			}, nil
		},

		// Workspace commands
		"workspace list": func() (cli.Command, error) {
			return &WorkspaceListCommand{
				Meta: *meta,
			}, nil
		},
		"workspace create": func() (cli.Command, error) {
			return &WorkspaceCreateCommand{
				Meta: *meta,
			}, nil
		},
		"workspace read": func() (cli.Command, error) {
			return &WorkspaceReadCommand{
				Meta: *meta,
			}, nil
		},
		"workspace update": func() (cli.Command, error) {
			return &WorkspaceUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"workspace delete": func() (cli.Command, error) {
			return &WorkspaceDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Run commands
		"run list": func() (cli.Command, error) {
			return &RunListCommand{
				Meta: *meta,
			}, nil
		},
		"run create": func() (cli.Command, error) {
			return &RunCreateCommand{
				Meta: *meta,
			}, nil
		},
		"run show": func() (cli.Command, error) {
			return &RunShowCommand{
				Meta: *meta,
			}, nil
		},
		"run apply": func() (cli.Command, error) {
			return &RunApplyCommand{
				Meta: *meta,
			}, nil
		},
		"run discard": func() (cli.Command, error) {
			return &RunDiscardCommand{
				Meta: *meta,
			}, nil
		},
		"run cancel": func() (cli.Command, error) {
			return &RunCancelCommand{
				Meta: *meta,
			}, nil
		},

		// Plan commands
		"plan read": func() (cli.Command, error) {
			return &PlanReadCommand{
				Meta: *meta,
			}, nil
		},
		"plan logs": func() (cli.Command, error) {
			return &PlanLogsCommand{
				Meta: *meta,
			}, nil
		},

		// Plan Export commands
		"planexport create": func() (cli.Command, error) {
			return &PlanExportCreateCommand{
				Meta: *meta,
			}, nil
		},
		"planexport read": func() (cli.Command, error) {
			return &PlanExportReadCommand{
				Meta: *meta,
			}, nil
		},
		"planexport download": func() (cli.Command, error) {
			return &PlanExportDownloadCommand{
				Meta: *meta,
			}, nil
		},
		"planexport delete": func() (cli.Command, error) {
			return &PlanExportDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Apply commands
		"apply read": func() (cli.Command, error) {
			return &ApplyReadCommand{
				Meta: *meta,
			}, nil
		},
		"apply logs": func() (cli.Command, error) {
			return &ApplyLogsCommand{
				Meta: *meta,
			}, nil
		},

		// Configuration Version commands
		"configversion list": func() (cli.Command, error) {
			return &ConfigVersionListCommand{
				Meta: *meta,
			}, nil
		},
		"configversion create": func() (cli.Command, error) {
			return &ConfigVersionCreateCommand{
				Meta: *meta,
			}, nil
		},
		"configversion read": func() (cli.Command, error) {
			return &ConfigVersionReadCommand{
				Meta: *meta,
			}, nil
		},
		"configversion upload": func() (cli.Command, error) {
			return &ConfigVersionUploadCommand{
				Meta: *meta,
			}, nil
		},

		// Organization commands
		"organization list": func() (cli.Command, error) {
			return &OrganizationListCommand{
				Meta: *meta,
			}, nil
		},
		"organization create": func() (cli.Command, error) {
			return &OrganizationCreateCommand{
				Meta: *meta,
			}, nil
		},
		"organization show": func() (cli.Command, error) {
			return &OrganizationShowCommand{
				Meta: *meta,
			}, nil
		},
		"organization update": func() (cli.Command, error) {
			return &OrganizationUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"organization delete": func() (cli.Command, error) {
			return &OrganizationDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Variable commands
		"variable list": func() (cli.Command, error) {
			return &VariableListCommand{
				Meta: *meta,
			}, nil
		},
		"variable create": func() (cli.Command, error) {
			return &VariableCreateCommand{
				Meta: *meta,
			}, nil
		},
		"variable update": func() (cli.Command, error) {
			return &VariableUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"variable delete": func() (cli.Command, error) {
			return &VariableDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Variable Set commands
		"variableset list": func() (cli.Command, error) {
			return &VariableSetListCommand{
				Meta: *meta,
			}, nil
		},
		"variableset create": func() (cli.Command, error) {
			return &VariableSetCreateCommand{
				Meta: *meta,
			}, nil
		},
		"variableset read": func() (cli.Command, error) {
			return &VariableSetReadCommand{
				Meta: *meta,
			}, nil
		},
		"variableset update": func() (cli.Command, error) {
			return &VariableSetUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"variableset delete": func() (cli.Command, error) {
			return &VariableSetDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"variableset apply": func() (cli.Command, error) {
			return &VariableSetApplyCommand{
				Meta: *meta,
			}, nil
		},

		// Variable Set Variable commands
		"variableset variable list": func() (cli.Command, error) {
			return &VariableSetVariableListCommand{
				Meta: *meta,
			}, nil
		},
		"variableset variable create": func() (cli.Command, error) {
			return &VariableSetVariableCreateCommand{
				Meta: *meta,
			}, nil
		},
		"variableset variable update": func() (cli.Command, error) {
			return &VariableSetVariableUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"variableset variable delete": func() (cli.Command, error) {
			return &VariableSetVariableDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Team commands
		"team list": func() (cli.Command, error) {
			return &TeamListCommand{
				Meta: *meta,
			}, nil
		},
		"team create": func() (cli.Command, error) {
			return &TeamCreateCommand{
				Meta: *meta,
			}, nil
		},
		"team show": func() (cli.Command, error) {
			return &TeamShowCommand{
				Meta: *meta,
			}, nil
		},
		"team delete": func() (cli.Command, error) {
			return &TeamDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"team add-member": func() (cli.Command, error) {
			return &TeamAddMemberCommand{
				Meta: *meta,
			}, nil
		},
		"team remove-member": func() (cli.Command, error) {
			return &TeamRemoveMemberCommand{
				Meta: *meta,
			}, nil
		},

		// Policy commands
		"policy list": func() (cli.Command, error) {
			return &PolicyListCommand{
				Meta: *meta,
			}, nil
		},
		"policy create": func() (cli.Command, error) {
			return &PolicyCreateCommand{
				Meta: *meta,
			}, nil
		},
		"policy read": func() (cli.Command, error) {
			return &PolicyReadCommand{
				Meta: *meta,
			}, nil
		},
		"policy update": func() (cli.Command, error) {
			return &PolicyUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"policy delete": func() (cli.Command, error) {
			return &PolicyDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Policy Set commands
		"policyset list": func() (cli.Command, error) {
			return &PolicySetListCommand{
				Meta: *meta,
			}, nil
		},
		"policyset create": func() (cli.Command, error) {
			return &PolicySetCreateCommand{
				Meta: *meta,
			}, nil
		},
		"policyset read": func() (cli.Command, error) {
			return &PolicySetReadCommand{
				Meta: *meta,
			}, nil
		},
		"policyset update": func() (cli.Command, error) {
			return &PolicySetUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"policyset delete": func() (cli.Command, error) {
			return &PolicySetDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"policyset add-policy": func() (cli.Command, error) {
			return &PolicySetAddPolicyCommand{
				Meta: *meta,
			}, nil
		},
		"policyset remove-policy": func() (cli.Command, error) {
			return &PolicySetRemovePolicyCommand{
				Meta: *meta,
			}, nil
		},

		// Policy Check commands (Sentinel policy check results)
		"policycheck list": func() (cli.Command, error) {
			return &PolicyCheckListCommand{
				Meta: *meta,
			}, nil
		},
		"policycheck read": func() (cli.Command, error) {
			return &PolicyCheckReadCommand{
				Meta: *meta,
			}, nil
		},
		"policycheck override": func() (cli.Command, error) {
			return &PolicyCheckOverrideCommand{
				Meta: *meta,
			}, nil
		},

		// Policy Evaluation commands (individual policy evaluation results in task stages)
		"policyevaluation list": func() (cli.Command, error) {
			return &PolicyEvaluationListCommand{
				Meta: *meta,
			}, nil
		},

		// Policy Set Outcome commands (policy set evaluation outcomes)
		"policyset outcome": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "policyset outcome",
				synopsis: "Manage policy set outcomes",
			}, nil
		},
		"policyset outcome list": func() (cli.Command, error) {
			return &PolicySetOutcomeListCommand{
				Meta: *meta,
			}, nil
		},
		"policyset outcome read": func() (cli.Command, error) {
			return &PolicySetOutcomeReadCommand{
				Meta: *meta,
			}, nil
		},
		// Policy Set Parameter commands (parameters for policy sets)
		"policyset parameter": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "policyset parameter",
				synopsis: "Manage policy set parameters",
			}, nil
		},
		"policyset parameter list": func() (cli.Command, error) {
			return &PolicySetParameterListCommand{
				Meta: *meta,
			}, nil
		},
		"policyset parameter create": func() (cli.Command, error) {
			return &PolicySetParameterCreateCommand{
				Meta: *meta,
			}, nil
		},
		"policyset parameter update": func() (cli.Command, error) {
			return &PolicySetParameterUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"policyset parameter delete": func() (cli.Command, error) {
			return &PolicySetParameterDeleteCommand{
				Meta: *meta,
			}, nil
		},
		// SSH Key commands
		"sshkey list": func() (cli.Command, error) {
			return &SSHKeyListCommand{
				Meta: *meta,
			}, nil
		},
		"sshkey create": func() (cli.Command, error) {
			return &SSHKeyCreateCommand{
				Meta: *meta,
			}, nil
		},
		"sshkey read": func() (cli.Command, error) {
			return &SSHKeyReadCommand{
				Meta: *meta,
			}, nil
		},
		"sshkey update": func() (cli.Command, error) {
			return &SSHKeyUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"sshkey delete": func() (cli.Command, error) {
			return &SSHKeyDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Project commands
		"project list": func() (cli.Command, error) {
			return &ProjectListCommand{
				Meta: *meta,
			}, nil
		},
		"project create": func() (cli.Command, error) {
			return &ProjectCreateCommand{
				Meta: *meta,
			}, nil
		},
		"project read": func() (cli.Command, error) {
			return &ProjectReadCommand{
				Meta: *meta,
			}, nil
		},
		"project update": func() (cli.Command, error) {
			return &ProjectUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"project delete": func() (cli.Command, error) {
			return &ProjectDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// State commands
		"state list": func() (cli.Command, error) {
			return &StateListCommand{
				Meta: *meta,
			}, nil
		},
		"state read": func() (cli.Command, error) {
			return &StateReadCommand{
				Meta: *meta,
			}, nil
		},
		"state outputs": func() (cli.Command, error) {
			return &StateOutputsCommand{
				Meta: *meta,
			}, nil
		},
		"state download": func() (cli.Command, error) {
			return &StateDownloadCommand{
				Meta: *meta,
			}, nil
		},

		// Notification commands
		"notification list": func() (cli.Command, error) {
			return &NotificationListCommand{
				Meta: *meta,
			}, nil
		},
		"notification create": func() (cli.Command, error) {
			return &NotificationCreateCommand{
				Meta: *meta,
			}, nil
		},
		"notification read": func() (cli.Command, error) {
			return &NotificationReadCommand{
				Meta: *meta,
			}, nil
		},
		"notification update": func() (cli.Command, error) {
			return &NotificationUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"notification delete": func() (cli.Command, error) {
			return &NotificationDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"notification verify": func() (cli.Command, error) {
			return &NotificationVerifyCommand{
				Meta: *meta,
			}, nil
		},

		// Run Task commands
		"runtask list": func() (cli.Command, error) {
			return &RunTaskListCommand{
				Meta: *meta,
			}, nil
		},
		"runtask create": func() (cli.Command, error) {
			return &RunTaskCreateCommand{
				Meta: *meta,
			}, nil
		},
		"runtask read": func() (cli.Command, error) {
			return &RunTaskReadCommand{
				Meta: *meta,
			}, nil
		},
		"runtask update": func() (cli.Command, error) {
			return &RunTaskUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"runtask delete": func() (cli.Command, error) {
			return &RunTaskDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"runtask attach": func() (cli.Command, error) {
			return &RunTaskAttachCommand{
				Meta: *meta,
			}, nil
		},
		"runtask detach": func() (cli.Command, error) {
			return &RunTaskDetachCommand{
				Meta: *meta,
			}, nil
		},

		// Run Trigger commands
		"runtrigger list": func() (cli.Command, error) {
			return &RunTriggerListCommand{
				Meta: *meta,
			}, nil
		},
		"runtrigger create": func() (cli.Command, error) {
			return &RunTriggerCreateCommand{
				Meta: *meta,
			}, nil
		},
		"runtrigger read": func() (cli.Command, error) {
			return &RunTriggerReadCommand{
				Meta: *meta,
			}, nil
		},
		"runtrigger delete": func() (cli.Command, error) {
			return &RunTriggerDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Agent Pool commands
		"agentpool list": func() (cli.Command, error) {
			return &AgentPoolListCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool create": func() (cli.Command, error) {
			return &AgentPoolCreateCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool read": func() (cli.Command, error) {
			return &AgentPoolReadCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool update": func() (cli.Command, error) {
			return &AgentPoolUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool delete": func() (cli.Command, error) {
			return &AgentPoolDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool token-create": func() (cli.Command, error) {
			return &AgentPoolTokenCreateCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool token-list": func() (cli.Command, error) {
			return &AgentPoolTokenListCommand{
				Meta: *meta,
			}, nil
		},
		"agentpool token-delete": func() (cli.Command, error) {
			return &AgentPoolTokenDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Agent commands (monitor self-hosted agent status)
		"agent list": func() (cli.Command, error) {
			return &AgentListCommand{
				Meta: *meta,
			}, nil
		},
		"agent read": func() (cli.Command, error) {
			return &AgentReadCommand{
				Meta: *meta,
			}, nil
		},

		// OAuth Client commands
		"oauthclient list": func() (cli.Command, error) {
			return &OAuthClientListCommand{
				Meta: *meta,
			}, nil
		},
		"oauthclient create": func() (cli.Command, error) {
			return &OAuthClientCreateCommand{
				Meta: *meta,
			}, nil
		},
		"oauthclient read": func() (cli.Command, error) {
			return &OAuthClientReadCommand{
				Meta: *meta,
			}, nil
		},
		"oauthclient update": func() (cli.Command, error) {
			return &OAuthClientUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"oauthclient delete": func() (cli.Command, error) {
			return &OAuthClientDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// OAuth Token commands
		"oauthtoken list": func() (cli.Command, error) {
			return &OAuthTokenListCommand{
				Meta: *meta,
			}, nil
		},
		"oauthtoken read": func() (cli.Command, error) {
			return &OAuthTokenReadCommand{
				Meta: *meta,
			}, nil
		},
		"oauthtoken update": func() (cli.Command, error) {
			return &OAuthTokenUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"oauthtoken delete": func() (cli.Command, error) {
			return &OAuthTokenDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Team Access commands (workspace permissions)
		"team access": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "team access",
				synopsis: "Manage team access",
			}, nil
		},
		"team access list": func() (cli.Command, error) {
			return &TeamAccessListCommand{
				Meta: *meta,
			}, nil
		},
		"team access create": func() (cli.Command, error) {
			return &TeamAccessCreateCommand{
				Meta: *meta,
			}, nil
		},
		"team access read": func() (cli.Command, error) {
			return &TeamAccessReadCommand{
				Meta: *meta,
			}, nil
		},
		"team access update": func() (cli.Command, error) {
			return &TeamAccessUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"team access delete": func() (cli.Command, error) {
			return &TeamAccessDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Project Team Access commands (project permissions)
		"project teamaccess": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "project teamaccess",
				synopsis: "Manage project team access",
			}, nil
		},
		"project teamaccess list": func() (cli.Command, error) {
			return &ProjectTeamAccessListCommand{
				Meta: *meta,
			}, nil
		},
		"project teamaccess create": func() (cli.Command, error) {
			return &ProjectTeamAccessCreateCommand{
				Meta: *meta,
			}, nil
		},
		"project teamaccess read": func() (cli.Command, error) {
			return &ProjectTeamAccessReadCommand{
				Meta: *meta,
			}, nil
		},
		"project teamaccess update": func() (cli.Command, error) {
			return &ProjectTeamAccessUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"project teamaccess delete": func() (cli.Command, error) {
			return &ProjectTeamAccessDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Audit Trail commands (compliance and security monitoring)
		"audittrail list": func() (cli.Command, error) {
			return &AuditTrailListCommand{
				Meta: *meta,
			}, nil
		},
		"audittrail read": func() (cli.Command, error) {
			return &AuditTrailReadCommand{
				Meta: *meta,
			}, nil
		},

		// Audit Trail Token commands (manage audit log streaming tokens)
		"audittrail token": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "audittrail token",
				synopsis: "Manage audit trail tokens",
			}, nil
		},
		"audittrail token list": func() (cli.Command, error) {
			return &AuditTrailTokenListCommand{
				Meta: *meta,
			}, nil
		},
		"audittrail token create": func() (cli.Command, error) {
			return &AuditTrailTokenCreateCommand{
				Meta: *meta,
			}, nil
		},
		"audittrail token read": func() (cli.Command, error) {
			return &AuditTrailTokenReadCommand{
				Meta: *meta,
			}, nil
		},
		"audittrail token delete": func() (cli.Command, error) {
			return &AuditTrailTokenDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Stack parent command
		"stack": func() (cli.Command, error) {
			return &StackCommand{
				Meta: *meta,
			}, nil
		},

		// Stack management commands (new hierarchical namespace)
		"stack list": func() (cli.Command, error) {
			return &StackListCommand{
				Meta: *meta,
			}, nil
		},
		"stack create": func() (cli.Command, error) {
			return &StackCreateCommand{
				Meta: *meta,
			}, nil
		},
		"stack read": func() (cli.Command, error) {
			return &StackReadCommand{
				Meta: *meta,
			}, nil
		},
		"stack update": func() (cli.Command, error) {
			return &StackUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"stack delete": func() (cli.Command, error) {
			return &StackDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Stack configuration commands (new hierarchical namespace)
		"stack configuration list": func() (cli.Command, error) {
			return &StackConfigurationListCommand{
				Meta: *meta,
			}, nil
		},
		"stack configuration create": func() (cli.Command, error) {
			return &StackConfigurationCreateCommand{
				Meta: *meta,
			}, nil
		},
		"stack configuration read": func() (cli.Command, error) {
			return &StackConfigurationReadCommand{
				Meta: *meta,
			}, nil
		},
		"stack configuration update": func() (cli.Command, error) {
			return &StackConfigurationUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"stack configuration delete": func() (cli.Command, error) {
			return &StackConfigurationDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Stack deployment commands (new hierarchical namespace)
		"stack deployment list": func() (cli.Command, error) {
			return &StackDeploymentListCommand{
				Meta: *meta,
			}, nil
		},
		"stack deployment create": func() (cli.Command, error) {
			return &StackDeploymentCreateCommand{
				Meta: *meta,
			}, nil
		},
		"stack deployment read": func() (cli.Command, error) {
			return &StackDeploymentReadCommand{
				Meta: *meta,
			}, nil
		},

		// Stack state commands (new hierarchical namespace)
		"stack state list": func() (cli.Command, error) {
			return &StackStateListCommand{
				Meta: *meta,
			}, nil
		},
		"stack state read": func() (cli.Command, error) {
			return &StackStateReadCommand{
				Meta: *meta,
			}, nil
		},

		// Registry parent command
		"registry": func() (cli.Command, error) {
			return &RegistryCommand{
				Meta: *meta,
			}, nil
		},

		// Registry Module commands (new hierarchical namespace)
		"registry module list": func() (cli.Command, error) {
			return &RegistryModuleListCommand{
				Meta: *meta,
			}, nil
		},
		"registry module create": func() (cli.Command, error) {
			return &RegistryModuleCreateCommand{
				Meta: *meta,
			}, nil
		},
		"registry module read": func() (cli.Command, error) {
			return &RegistryModuleReadCommand{
				Meta: *meta,
			}, nil
		},
		"registry module delete": func() (cli.Command, error) {
			return &RegistryModuleDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"registry module version create": func() (cli.Command, error) {
			return &RegistryModuleCreateVersionCommand{
				Meta: *meta,
			}, nil
		},
		"registry module version delete": func() (cli.Command, error) {
			return &RegistryModuleDeleteVersionCommand{
				Meta: *meta,
			}, nil
		},

		// Registry Provider commands (new hierarchical namespace)
		"registry provider list": func() (cli.Command, error) {
			return &RegistryProviderListCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider create": func() (cli.Command, error) {
			return &RegistryProviderCreateCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider read": func() (cli.Command, error) {
			return &RegistryProviderReadCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider delete": func() (cli.Command, error) {
			return &RegistryProviderDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider version create": func() (cli.Command, error) {
			return &RegistryProviderVersionCreateCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider version read": func() (cli.Command, error) {
			return &RegistryProviderVersionReadCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider version delete": func() (cli.Command, error) {
			return &RegistryProviderVersionDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider platform create": func() (cli.Command, error) {
			return &RegistryProviderPlatformCreateCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider platform read": func() (cli.Command, error) {
			return &RegistryProviderPlatformReadCommand{
				Meta: *meta,
			}, nil
		},
		"registry provider platform delete": func() (cli.Command, error) {
			return &RegistryProviderPlatformDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Public Registry commands (query public Terraform registry)
		"publicregistry": func() (cli.Command, error) {
			return &PublicRegistryCommand{
				Meta: *meta,
			}, nil
		},
		"publicregistry provider": func() (cli.Command, error) {
			return &PublicRegistryProviderCommand{
				Meta: *meta,
			}, nil
		},
		"publicregistry provider versions": func() (cli.Command, error) {
			return &PublicRegistryProviderVersionsCommand{
				Meta: *meta,
			}, nil
		},
		"publicregistry module": func() (cli.Command, error) {
			return &PublicRegistryModuleCommand{
				Meta: *meta,
			}, nil
		},
		"publicregistry policy": func() (cli.Command, error) {
			return &PublicRegistryPolicyCommand{
				Meta: *meta,
			}, nil
		},
		"publicregistry policy list": func() (cli.Command, error) {
			return &PublicRegistryPolicyListCommand{
				Meta: *meta,
			}, nil
		},

		// GPG Key commands (manage GPG keys for provider signing)
		"gpgkey list": func() (cli.Command, error) {
			return &GPGKeyListCommand{
				Meta: *meta,
			}, nil
		},
		"gpgkey create": func() (cli.Command, error) {
			return &GPGKeyCreateCommand{
				Meta: *meta,
			}, nil
		},
		"gpgkey read": func() (cli.Command, error) {
			return &GPGKeyReadCommand{
				Meta: *meta,
			}, nil
		},
		"gpgkey update": func() (cli.Command, error) {
			return &GPGKeyUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"gpgkey delete": func() (cli.Command, error) {
			return &GPGKeyDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Cost estimate commands
		"costestimate read": func() (cli.Command, error) {
			return &CostEstimateReadCommand{
				Meta: *meta,
			}, nil
		},

		// Feature set commands
		"featureset list": func() (cli.Command, error) {
			return &FeatureSetListCommand{
				Meta: *meta,
			}, nil
		},

		// GitHub app commands
		"githubapp list": func() (cli.Command, error) {
			return &GitHubAppListCommand{
				Meta: *meta,
			}, nil
		},
		"githubapp read": func() (cli.Command, error) {
			return &GitHubAppReadCommand{
				Meta: *meta,
			}, nil
		},

		// IP range commands
		"iprange list": func() (cli.Command, error) {
			return &IPRangeReadCommand{
				Meta: *meta,
			}, nil
		},

		// No-code provisioning commands
		"nocode list": func() (cli.Command, error) {
			return &NoCodeListCommand{
				Meta: *meta,
			}, nil
		},
		"nocode create": func() (cli.Command, error) {
			return &NoCodeCreateCommand{
				Meta: *meta,
			}, nil
		},
		"nocode read": func() (cli.Command, error) {
			return &NoCodeReadCommand{
				Meta: *meta,
			}, nil
		},
		"nocode update": func() (cli.Command, error) {
			return &NoCodeUpdateCommand{
				Meta: *meta,
			}, nil
		},

		// Stability policy commands
		"stabilitypolicy read": func() (cli.Command, error) {
			return &StabilityPolicyReadCommand{
				Meta: *meta,
			}, nil
		},

		// Subscription commands
		"subscription list": func() (cli.Command, error) {
			return &SubscriptionListCommand{
				Meta: *meta,
			}, nil
		},
		"subscription read": func() (cli.Command, error) {
			return &SubscriptionReadCommand{
				Meta: *meta,
			}, nil
		},

		// User commands
		"user read": func() (cli.Command, error) {
			return &UserReadCommand{
				Meta: *meta,
			}, nil
		},

		// User Token commands (manage user-level API tokens)
		"user token": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "user token",
				synopsis: "Manage user tokens",
			}, nil
		},
		"user token list": func() (cli.Command, error) {
			return &UserTokenListCommand{
				Meta: *meta,
			}, nil
		},
		"user token create": func() (cli.Command, error) {
			return &UserTokenCreateCommand{
				Meta: *meta,
			}, nil
		},
		"user token read": func() (cli.Command, error) {
			return &UserTokenReadCommand{
				Meta: *meta,
			}, nil
		},
		"user token delete": func() (cli.Command, error) {
			return &UserTokenDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Team Token commands (manage team-level API tokens)
		"team token": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "team token",
				synopsis: "Manage team tokens",
			}, nil
		},
		"team token list": func() (cli.Command, error) {
			return &TeamTokenListCommand{
				Meta: *meta,
			}, nil
		},
		"team token create": func() (cli.Command, error) {
			return &TeamTokenCreateCommand{
				Meta: *meta,
			}, nil
		},
		"team token read": func() (cli.Command, error) {
			return &TeamTokenReadCommand{
				Meta: *meta,
			}, nil
		},
		"team token delete": func() (cli.Command, error) {
			return &TeamTokenDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Organization Membership commands (manage organization membership)
		"organization membership": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "organization membership",
				synopsis: "Manage organization memberships",
			}, nil
		},
		"organization membership list": func() (cli.Command, error) {
			return &OrganizationMembershipListCommand{
				Meta: *meta,
			}, nil
		},
		"organization membership create": func() (cli.Command, error) {
			return &OrganizationMembershipCreateCommand{
				Meta: *meta,
			}, nil
		},
		"organization membership read": func() (cli.Command, error) {
			return &OrganizationMembershipReadCommand{
				Meta: *meta,
			}, nil
		},
		"organization membership delete": func() (cli.Command, error) {
			return &OrganizationMembershipDeleteCommand{
				Meta: *meta,
			}, nil
		},
		// Organization Member commands (detailed member operations)
		"organization member": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "organization member",
				synopsis: "Manage organization members",
			}, nil
		},
		"organization member read": func() (cli.Command, error) {
			return &OrganizationMemberReadCommand{
				Meta: *meta,
			}, nil
		},
		// Organization Token commands
		"organization token": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "organization token",
				synopsis: "Manage organization tokens",
			}, nil
		},
		"organization token list": func() (cli.Command, error) {
			return &OrganizationTokenListCommand{
				Meta: *meta,
			}, nil
		},
		"organization token create": func() (cli.Command, error) {
			return &OrganizationTokenCreateCommand{
				Meta: *meta,
			}, nil
		},
		"organization token read": func() (cli.Command, error) {
			return &OrganizationTokenReadCommand{
				Meta: *meta,
			}, nil
		},
		"organization token delete": func() (cli.Command, error) {
			return &OrganizationTokenDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Organization Tag commands (tag resources in an organization)
		"organization tag": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "organization tag",
				synopsis: "Manage organization tags",
			}, nil
		},
		"organization tag list": func() (cli.Command, error) {
			return &OrganizationTagListCommand{
				Meta: *meta,
			}, nil
		},
		"organization tag create": func() (cli.Command, error) {
			return &OrganizationTagCreateCommand{
				Meta: *meta,
			}, nil
		},
		"organization tag delete": func() (cli.Command, error) {
			return &OrganizationTagDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Workspace Tag commands (apply organization tags to workspaces)
		"workspace tag": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "workspace tag",
				synopsis: "Manage workspace tags",
			}, nil
		},
		"workspace tag list": func() (cli.Command, error) {
			return &WorkspaceTagListCommand{
				Meta: *meta,
			}, nil
		},
		"workspace tag add": func() (cli.Command, error) {
			return &WorkspaceTagAddCommand{
				Meta: *meta,
			}, nil
		},
		"workspace tag remove": func() (cli.Command, error) {
			return &WorkspaceTagRemoveCommand{
				Meta: *meta,
			}, nil
		},

		// Reserved Tag Key commands (reserved tag key management)
		"reservedtagkey list": func() (cli.Command, error) {
			return &ReservedTagKeyListCommand{
				Meta: *meta,
			}, nil
		},
		"reservedtagkey create": func() (cli.Command, error) {
			return &ReservedTagKeyCreateCommand{
				Meta: *meta,
			}, nil
		},
		"reservedtagkey update": func() (cli.Command, error) {
			return &ReservedTagKeyUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"reservedtagkey delete": func() (cli.Command, error) {
			return &ReservedTagKeyDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Comment commands (run discussion/comments)
		"comment list": func() (cli.Command, error) {
			return &CommentListCommand{
				Meta: *meta,
			}, nil
		},
		"comment create": func() (cli.Command, error) {
			return &CommentCreateCommand{
				Meta: *meta,
			}, nil
		},
		"comment read": func() (cli.Command, error) {
			return &CommentReadCommand{
				Meta: *meta,
			}, nil
		},

		// AWS OIDC Configuration commands (dynamic AWS credentials)
		"awsoidc create": func() (cli.Command, error) {
			return &AWSoidcCreateCommand{
				Meta: *meta,
			}, nil
		},
		"awsoidc read": func() (cli.Command, error) {
			return &AWSoidcReadCommand{
				Meta: *meta,
			}, nil
		},
		"awsoidc update": func() (cli.Command, error) {
			return &AWSoidcUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"awsoidc delete": func() (cli.Command, error) {
			return &AWSoidcDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Azure OIDC Configuration commands (dynamic Azure credentials)
		"azureoidc create": func() (cli.Command, error) {
			return &AzureoidcCreateCommand{
				Meta: *meta,
			}, nil
		},
		"azureoidc read": func() (cli.Command, error) {
			return &AzureoidcReadCommand{
				Meta: *meta,
			}, nil
		},
		"azureoidc update": func() (cli.Command, error) {
			return &AzureoidcUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"azureoidc delete": func() (cli.Command, error) {
			return &AzureoidcDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// GCP OIDC Configuration commands (dynamic GCP credentials)
		"gcpoidc create": func() (cli.Command, error) {
			return &GCPoidcCreateCommand{
				Meta: *meta,
			}, nil
		},
		"gcpoidc read": func() (cli.Command, error) {
			return &GCPoidcReadCommand{
				Meta: *meta,
			}, nil
		},
		"gcpoidc update": func() (cli.Command, error) {
			return &GCPoidcUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"gcpoidc delete": func() (cli.Command, error) {
			return &GCPoidcDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Vault OIDC Configuration commands (dynamic Vault credentials)
		"vaultoidc create": func() (cli.Command, error) {
			return &VaultoidcCreateCommand{
				Meta: *meta,
			}, nil
		},
		"vaultoidc read": func() (cli.Command, error) {
			return &VaultoidcReadCommand{
				Meta: *meta,
			}, nil
		},
		"vaultoidc update": func() (cli.Command, error) {
			return &VaultoidcUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"vaultoidc delete": func() (cli.Command, error) {
			return &VaultoidcDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// Workspace Resource commands (view managed infrastructure)
		"workspace resource": func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     "workspace resource",
				synopsis: "Manage workspace resources",
			}, nil
		},
		"workspace resource list": func() (cli.Command, error) {
			return &WorkspaceResourceListCommand{
				Meta: *meta,
			}, nil
		},
		"workspace resource read": func() (cli.Command, error) {
			return &WorkspaceResourceReadCommand{
				Meta: *meta,
			}, nil
		},

		// Change Request commands (workspace to-do tracking and compliance management)
		"changerequest list": func() (cli.Command, error) {
			return &ChangeRequestListCommand{
				Meta: *meta,
			}, nil
		},
		"changerequest create": func() (cli.Command, error) {
			return &ChangeRequestCreateCommand{
				Meta: *meta,
			}, nil
		},
		"changerequest read": func() (cli.Command, error) {
			return &ChangeRequestReadCommand{
				Meta: *meta,
			}, nil
		},
		"changerequest update": func() (cli.Command, error) {
			return &ChangeRequestUpdateCommand{
				Meta: *meta,
			}, nil
		},

		// Assessment Result commands (health assessments, drift detection, continuous validation)
		"assessmentresult list": func() (cli.Command, error) {
			return &AssessmentResultListCommand{
				Meta: *meta,
			}, nil
		},
		"assessmentresult read": func() (cli.Command, error) {
			return &AssessmentResultReadCommand{
				Meta: *meta,
			}, nil
		},

		// Query commands (search across organization)
		"queryrun list": func() (cli.Command, error) {
			return &QueryRunListCommand{
				Meta: *meta,
			}, nil
		},
		"queryworkspace list": func() (cli.Command, error) {
			return &QueryWorkspaceListCommand{
				Meta: *meta,
			}, nil
		},

		// Explorer API commands (query resources across organization)
		"explorer query": func() (cli.Command, error) {
			return &ExplorerQueryCommand{
				Meta: *meta,
			}, nil
		},

		// VCS Event commands (VCS integration debugging and monitoring)
		"vcsevent list": func() (cli.Command, error) {
			return &VCSEventListCommand{
				Meta: *meta,
			}, nil
		},
		"vcsevent read": func() (cli.Command, error) {
			return &VCSEventReadCommand{
				Meta: *meta,
			}, nil
		},

		// HYOK (Hold Your Own Key) Configuration commands (customer-managed encryption keys)
		"hyok list": func() (cli.Command, error) {
			return &HYOKListCommand{
				Meta: *meta,
			}, nil
		},
		"hyok create": func() (cli.Command, error) {
			return &HYOKCreateCommand{
				Meta: *meta,
			}, nil
		},
		"hyok read": func() (cli.Command, error) {
			return &HYOKReadCommand{
				Meta: *meta,
			}, nil
		},
		"hyok update": func() (cli.Command, error) {
			return &HYOKUpdateCommand{
				Meta: *meta,
			}, nil
		},
		"hyok delete": func() (cli.Command, error) {
			return &HYOKDeleteCommand{
				Meta: *meta,
			}, nil
		},

		// HYOK Customer Key Version commands (manage key version lifecycle)
		"hyokkey create": func() (cli.Command, error) {
			return &HYOKKeyCreateCommand{
				Meta: *meta,
			}, nil
		},
		"hyokkey read": func() (cli.Command, error) {
			return &HYOKKeyReadCommand{
				Meta: *meta,
			}, nil
		},
		"hyokkey delete": func() (cli.Command, error) {
			return &HYOKKeyDeleteCommand{
				Meta: *meta,
			}, nil
		},
	}

	namespaceSynopses := map[string]string{
		"costestimate":     "Manage cost estimates",
		"featureset":       "Manage feature sets",
		"githubapp":        "Manage GitHub app installations",
		"iprange":          "View Terraform IP ranges",
		"nocode":           "Manage no-code provisioning",
		"account":          "Manage accounts",
		"agent":            "Manage agents",
		"apply":            "Manage applies",
		"assessmentresult": "Manage assessment results",
		"audittrail":       "Manage audit trail entries",
		"audittrail token": "Manage audit trail tokens",
		"awsoidc":          "Manage AWS OIDC integration",
		"azureoidc":        "Manage Azure OIDC integration",
		"changerequest":    "Manage change requests",
		"comment":          "Manage run comments",
		"configversion":    "Manage workspace configuration versions",
		"explorer":         "Query Terraform Cloud",
		"gcpoidc":          "Manage GCP OIDC integration",
		"gpgkey":           "Manage GPG keys",
		"hyok":             "Manage Hold Your Own Key settings",
		"hyokkey":          "Manage Hold Your Own Key versions",
		"notification":     "Manage notifications",
		"oauthclient":      "Manage OAuth clients",
		"oauthtoken":       "Manage OAuth tokens",
		"organization":     "Manage the current organization",
		"organization tag": "Manage organization tags",
		"plan":             "Manage Terraform plans",
		"planexport":       "Manage plan exports",
		"policy":           "Manage policies",
		"policycheck":      "Manage policy checks",
		"policyevaluation": "Manage policy evaluations",
		"policyset":        "Manage policy sets",
		"project":          "Manage projects",
		"queryrun":         "Search runs",
		"queryworkspace":   "Search workspaces",
		"stabilitypolicy":  "Read stability policy",
		"subscription":     "Manage subscriptions",
		"reservedtagkey":   "Manage reserved tag keys",
		"run":              "Manage Terraform runs",
		"runtask":          "Manage run tasks",
		"runtrigger":       "Manage run triggers",
		"sshkey":           "Manage SSH keys",
		"state":            "Manage Terraform states",
		"team":             "Manage teams",
		"team token":       "Manage team tokens",
		"user token":       "Manage user tokens",
		"variable":         "Manage workspace variables",
		"variableset":      "Manage variable sets",
		"vaultoidc":        "Manage Vault OIDC integration",
		"vcsevent":         "Manage VCS events",
		"workspace":        "Manage workspaces",
	}

	for commandName := range commands {
		parts := strings.Split(commandName, " ")
		if len(parts) <= 1 {
			continue
		}

		namespace := parts[0]
		if _, exists := commands[namespace]; exists {
			continue
		}

		synopsis := "Manage " + namespace
		if customSynopsis, ok := namespaceSynopses[namespace]; ok {
			synopsis = customSynopsis
		}

		ns := namespace
		syn := synopsis
		commands[ns] = func() (cli.Command, error) {
			return &NamespaceCommand{
				Meta:     *meta,
				name:     ns,
				synopsis: syn,
			}, nil
		}
	}

	return commands
}
