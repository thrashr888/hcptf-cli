package command

import (
	"context"
	"io"

	tfe "github.com/hashicorp/go-tfe"
)

type mockWorkspaceCreateService struct {
	response    *tfe.Workspace
	err         error
	lastOrg     string
	lastOptions tfe.WorkspaceCreateOptions
}

func (m *mockWorkspaceCreateService) Create(_ context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

type mockWorkspaceUpdateService struct {
	response    *tfe.Workspace
	err         error
	lastOrg     string
	lastName    string
	lastOptions tfe.WorkspaceUpdateOptions
}

func (m *mockWorkspaceUpdateService) Update(_ context.Context, organization, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	m.lastOrg = organization
	m.lastName = workspace
	m.lastOptions = options
	return m.response, m.err
}

type mockWorkspaceDeleteService struct {
	err      error
	lastOrg  string
	lastName string
}

func (m *mockWorkspaceDeleteService) Delete(_ context.Context, organization, workspace string) error {
	m.lastOrg = organization
	m.lastName = workspace
	return m.err
}

type mockRunCreateService struct {
	response    *tfe.Run
	err         error
	lastOptions tfe.RunCreateOptions
}

func (m *mockRunCreateService) Create(_ context.Context, options tfe.RunCreateOptions) (*tfe.Run, error) {
	m.lastOptions = options
	return m.response, m.err
}

type mockRunApplyService struct {
	err         error
	lastRun     string
	lastOptions tfe.RunApplyOptions
}

func (m *mockRunApplyService) Apply(_ context.Context, runID string, options tfe.RunApplyOptions) error {
	m.lastRun = runID
	m.lastOptions = options
	return m.err
}

type mockPlanService struct {
	response *tfe.Plan
	err      error
	lastID   string
}

func (m *mockPlanService) Read(_ context.Context, planID string) (*tfe.Plan, error) {
	m.lastID = planID
	return m.response, m.err
}

type mockApplyService struct {
	response *tfe.Apply
	err      error
	lastID   string
}

func (m *mockApplyService) Read(_ context.Context, applyID string) (*tfe.Apply, error) {
	m.lastID = applyID
	return m.response, m.err
}

type mockConfigVersionListService struct {
	response      *tfe.ConfigurationVersionList
	err           error
	lastWorkspace string
	lastOptions   *tfe.ConfigurationVersionListOptions
}

func (m *mockConfigVersionListService) List(_ context.Context, workspaceID string, options *tfe.ConfigurationVersionListOptions) (*tfe.ConfigurationVersionList, error) {
	m.lastWorkspace = workspaceID
	if options != nil {
		copy := *options
		m.lastOptions = &copy
	}
	return m.response, m.err
}

type mockConfigVersionReadService struct {
	response *tfe.ConfigurationVersion
	err      error
	lastID   string
}

func (m *mockConfigVersionReadService) Read(_ context.Context, configurationID string) (*tfe.ConfigurationVersion, error) {
	m.lastID = configurationID
	return m.response, m.err
}

func (m *mockConfigVersionReadService) ReadWithOptions(ctx context.Context, configurationID string, options *tfe.ConfigurationVersionReadOptions) (*tfe.ConfigurationVersion, error) {
	return m.Read(ctx, configurationID)
}

type mockPlanLogService struct {
	reader io.Reader
	err    error
	lastID string
}

func (m *mockPlanLogService) Logs(_ context.Context, planID string) (io.Reader, error) {
	m.lastID = planID
	return m.reader, m.err
}

type mockApplyLogService struct {
	reader io.Reader
	err    error
	lastID string
}

func (m *mockApplyLogService) Logs(_ context.Context, applyID string) (io.Reader, error) {
	m.lastID = applyID
	return m.reader, m.err
}

type mockRunCancelService struct {
	cancelErr     error
	forceErr      error
	lastCancelRun string
	lastForceRun  string
	lastCancelOpt tfe.RunCancelOptions
	lastForceOpt  tfe.RunForceCancelOptions
}

func (m *mockRunCancelService) Cancel(_ context.Context, runID string, options tfe.RunCancelOptions) error {
	m.lastCancelRun = runID
	m.lastCancelOpt = options
	return m.cancelErr
}

func (m *mockRunCancelService) ForceCancel(_ context.Context, runID string, options tfe.RunForceCancelOptions) error {
	m.lastForceRun = runID
	m.lastForceOpt = options
	return m.forceErr
}

type mockRunDiscardService struct {
	err        error
	lastRun    string
	lastOption tfe.RunDiscardOptions
}

func (m *mockRunDiscardService) Discard(_ context.Context, runID string, options tfe.RunDiscardOptions) error {
	m.lastRun = runID
	m.lastOption = options
	return m.err
}

type mockRunReadService struct {
	response *tfe.Run
	err      error
	lastRun  string
}

func (m *mockRunReadService) Read(_ context.Context, runID string) (*tfe.Run, error) {
	m.lastRun = runID
	return m.response, m.err
}

type mockVariableCreateService struct {
	response      *tfe.Variable
	err           error
	lastWorkspace string
	lastOptions   tfe.VariableCreateOptions
}

func (m *mockVariableCreateService) Create(_ context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	m.lastWorkspace = workspaceID
	m.lastOptions = options
	return m.response, m.err
}

type mockVariableUpdateService struct {
	response      *tfe.Variable
	err           error
	lastWorkspace string
	lastID        string
	lastOptions   tfe.VariableUpdateOptions
}

func (m *mockVariableUpdateService) Update(_ context.Context, workspaceID, variableID string, options tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	m.lastWorkspace = workspaceID
	m.lastID = variableID
	m.lastOptions = options
	return m.response, m.err
}

type mockVariableDeleteService struct {
	err           error
	lastWorkspace string
	lastID        string
}

func (m *mockVariableDeleteService) Delete(_ context.Context, workspaceID, variableID string) error {
	m.lastWorkspace = workspaceID
	m.lastID = variableID
	return m.err
}

type mockGPGKeyCreateService struct {
	response     *tfe.GPGKey
	err          error
	lastRegistry tfe.RegistryName
	lastOptions  tfe.GPGKeyCreateOptions
}

func (m *mockGPGKeyCreateService) Create(_ context.Context, registryName tfe.RegistryName, options tfe.GPGKeyCreateOptions) (*tfe.GPGKey, error) {
	m.lastRegistry = registryName
	m.lastOptions = options
	return m.response, m.err
}

type mockGPGKeyListService struct {
	response    *tfe.GPGKeyList
	err         error
	lastOptions tfe.GPGKeyListOptions
}

func (m *mockGPGKeyListService) ListPrivate(_ context.Context, options tfe.GPGKeyListOptions) (*tfe.GPGKeyList, error) {
	m.lastOptions = options
	return m.response, m.err
}

type mockGPGKeyReadService struct {
	response *tfe.GPGKey
	err      error
	lastID   tfe.GPGKeyID
}

func (m *mockGPGKeyReadService) Read(_ context.Context, keyID tfe.GPGKeyID) (*tfe.GPGKey, error) {
	m.lastID = keyID
	return m.response, m.err
}

type mockGPGKeyUpdateService struct {
	response    *tfe.GPGKey
	err         error
	lastID      tfe.GPGKeyID
	lastOptions tfe.GPGKeyUpdateOptions
}

func (m *mockGPGKeyUpdateService) Update(_ context.Context, keyID tfe.GPGKeyID, options tfe.GPGKeyUpdateOptions) (*tfe.GPGKey, error) {
	m.lastID = keyID
	m.lastOptions = options
	return m.response, m.err
}

type mockGPGKeyDeleteService struct {
	err    error
	lastID tfe.GPGKeyID
}

func (m *mockGPGKeyDeleteService) Delete(_ context.Context, keyID tfe.GPGKeyID) error {
	m.lastID = keyID
	return m.err
}

type mockRegistryProviderListService struct {
	response *tfe.RegistryProviderList
	err      error
	lastOrg  string
	lastOpts *tfe.RegistryProviderListOptions
}

func (m *mockRegistryProviderListService) List(_ context.Context, organization string, options *tfe.RegistryProviderListOptions) (*tfe.RegistryProviderList, error) {
	m.lastOrg = organization
	if options != nil {
		copy := *options
		m.lastOpts = &copy
	}
	return m.response, m.err
}

type mockRegistryProviderCreateService struct {
	response *tfe.RegistryProvider
	err      error
	lastOrg  string
	lastOpts tfe.RegistryProviderCreateOptions
}

func (m *mockRegistryProviderCreateService) Create(_ context.Context, organization string, options tfe.RegistryProviderCreateOptions) (*tfe.RegistryProvider, error) {
	m.lastOrg = organization
	m.lastOpts = options
	return m.response, m.err
}

type mockRegistryProviderReadService struct {
	response *tfe.RegistryProvider
	err      error
	lastID   tfe.RegistryProviderID
}

func (m *mockRegistryProviderReadService) Read(_ context.Context, providerID tfe.RegistryProviderID, _ *tfe.RegistryProviderReadOptions) (*tfe.RegistryProvider, error) {
	m.lastID = providerID
	return m.response, m.err
}

type mockRegistryProviderDeleteService struct {
	err    error
	lastID tfe.RegistryProviderID
}

func (m *mockRegistryProviderDeleteService) Delete(_ context.Context, providerID tfe.RegistryProviderID) error {
	m.lastID = providerID
	return m.err
}
type mockOrganizationListService struct {
	response    *tfe.OrganizationList
	err         error
	lastOptions *tfe.OrganizationListOptions
}

func (m *mockOrganizationListService) List(_ context.Context, options *tfe.OrganizationListOptions) (*tfe.OrganizationList, error) {
	m.lastOptions = options
	return m.response, m.err
}

type mockTeamListService struct {
	response *tfe.TeamList
	err      error
	lastOrg  string
}

func (m *mockTeamListService) List(_ context.Context, organization string, _ *tfe.TeamListOptions) (*tfe.TeamList, error) {
	m.lastOrg = organization
	return m.response, m.err
}

type mockProjectListService struct {
	response *tfe.ProjectList
	err      error
	lastOrg  string
}

func (m *mockProjectListService) List(_ context.Context, organization string, _ *tfe.ProjectListOptions) (*tfe.ProjectList, error) {
	m.lastOrg = organization
	return m.response, m.err
}

type mockPolicyListService struct {
	response *tfe.PolicyList
	err      error
	lastOrg  string
}

func (m *mockPolicyListService) List(_ context.Context, organization string, _ *tfe.PolicyListOptions) (*tfe.PolicyList, error) {
	m.lastOrg = organization
	return m.response, m.err
}

type mockPolicySetListService struct {
	response *tfe.PolicySetList
	err      error
	lastOrg  string
}

func (m *mockPolicySetListService) List(_ context.Context, organization string, _ *tfe.PolicySetListOptions) (*tfe.PolicySetList, error) {
	m.lastOrg = organization
	return m.response, m.err
}

type mockOrganizationReadService struct {
	response *tfe.Organization
	err      error
	lastName string
}

func (m *mockOrganizationReadService) Read(_ context.Context, organization string) (*tfe.Organization, error) {
	m.lastName = organization
	return m.response, m.err
}

type mockTeamReadService struct {
	response *tfe.Team
	err      error
	lastName string
}

func (m *mockTeamReadService) Read(_ context.Context, teamName string) (*tfe.Team, error) {
	m.lastName = teamName
	return m.response, m.err
}

type mockProjectReadService struct {
	response *tfe.Project
	err      error
	lastID   string
}

func (m *mockProjectReadService) Read(_ context.Context, projectID string) (*tfe.Project, error) {
	m.lastID = projectID
	return m.response, m.err
}

type mockPolicyReadService struct {
	response *tfe.Policy
	err      error
	lastID   string
}

func (m *mockPolicyReadService) Read(_ context.Context, policyID string) (*tfe.Policy, error) {
	m.lastID = policyID
	return m.response, m.err
}

type mockPolicyDownloadService struct {
	content []byte
	err     error
	lastID  string
}

func (m *mockPolicyDownloadService) Download(_ context.Context, policyID string) ([]byte, error) {
	m.lastID = policyID
	return m.content, m.err
}

type mockPolicySetReadService struct {
	response *tfe.PolicySet
	err      error
	lastID   string
}

func (m *mockPolicySetReadService) Read(_ context.Context, policySetID string) (*tfe.PolicySet, error) {
	m.lastID = policySetID
	return m.response, m.err
}

type mockVariableSetReadService struct {
	response *tfe.VariableSet
	err      error
	lastID   string
}

func (m *mockVariableSetReadService) Read(_ context.Context, variableSetID string, _ *tfe.VariableSetReadOptions) (*tfe.VariableSet, error) {
	m.lastID = variableSetID
	return m.response, m.err
}

type mockVariableSetListService struct {
	response *tfe.VariableSetList
	err      error
	lastOrg  string
	lastOpts *tfe.VariableSetListOptions
}

func (m *mockVariableSetListService) List(_ context.Context, organization string, options *tfe.VariableSetListOptions) (*tfe.VariableSetList, error) {
	m.lastOrg = organization
	if options != nil {
		copy := *options
		m.lastOpts = &copy
	}
	return m.response, m.err
}

type mockVariableSetCreateService struct {
	response    *tfe.VariableSet
	err         error
	lastOrg     string
	lastOptions *tfe.VariableSetCreateOptions
}

func (m *mockVariableSetCreateService) Create(_ context.Context, organization string, options *tfe.VariableSetCreateOptions) (*tfe.VariableSet, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

type mockStateVersionReadService struct {
	response *tfe.StateVersion
	err      error
	lastID   string
}

func (m *mockStateVersionReadService) ReadCurrent(_ context.Context, workspaceID string) (*tfe.StateVersion, error) {
	m.lastID = workspaceID
	return m.response, m.err
}

type mockSSHKeyCreateService struct {
	response    *tfe.SSHKey
	err         error
	lastOrg     string
	lastOptions tfe.SSHKeyCreateOptions
}

func (m *mockSSHKeyCreateService) Create(_ context.Context, organization string, options tfe.SSHKeyCreateOptions) (*tfe.SSHKey, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

type mockAccountReadService struct {
	response *tfe.User
	err      error
}

func (m *mockAccountReadService) ReadCurrent(_ context.Context) (*tfe.User, error) {
	return m.response, m.err
}

type mockRunTaskCreateService struct {
	response    *tfe.RunTask
	err         error
	lastOrg     string
	lastOptions tfe.RunTaskCreateOptions
}

func (m *mockRunTaskCreateService) Create(_ context.Context, organization string, options tfe.RunTaskCreateOptions) (*tfe.RunTask, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

type mockRunTaskUpdateService struct {
	response    *tfe.RunTask
	err         error
	lastID      string
	lastOptions tfe.RunTaskUpdateOptions
}

func (m *mockRunTaskUpdateService) Update(_ context.Context, runTaskID string, options tfe.RunTaskUpdateOptions) (*tfe.RunTask, error) {
	m.lastID = runTaskID
	m.lastOptions = options
	return m.response, m.err
}

type mockRegistryProviderVersionCreateService struct {
	response    *tfe.RegistryProviderVersion
	err         error
	lastID      tfe.RegistryProviderID
	lastOptions tfe.RegistryProviderVersionCreateOptions
}

func (m *mockRegistryProviderVersionCreateService) Create(_ context.Context, providerID tfe.RegistryProviderID, options tfe.RegistryProviderVersionCreateOptions) (*tfe.RegistryProviderVersion, error) {
	m.lastID = providerID
	m.lastOptions = options
	return m.response, m.err
}

type mockRegistryProviderVersionDeleteService struct {
	err    error
	lastID tfe.RegistryProviderVersionID
}

func (m *mockRegistryProviderVersionDeleteService) Delete(_ context.Context, versionID tfe.RegistryProviderVersionID) error {
	m.lastID = versionID
	return m.err
}

type mockCommentListService struct {
	response *tfe.CommentList
	err      error
	lastID   string
}

func (m *mockCommentListService) List(_ context.Context, runID string) (*tfe.CommentList, error) {
	m.lastID = runID
	return m.response, m.err
}

type mockAgentPoolCreateService struct {
	response    *tfe.AgentPool
	err         error
	lastOrg     string
	lastOptions tfe.AgentPoolCreateOptions
}

func (m *mockAgentPoolCreateService) Create(_ context.Context, organization string, options tfe.AgentPoolCreateOptions) (*tfe.AgentPool, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

type mockUserReadService struct {
	response *tfe.User
	err      error
}

func (m *mockUserReadService) ReadCurrent(_ context.Context) (*tfe.User, error) {
	return m.response, m.err
}

type mockUserTokenListService struct {
	response *tfe.UserTokenList
	err      error
	lastID   string
}

func (m *mockUserTokenListService) List(_ context.Context, userID string) (*tfe.UserTokenList, error) {
	m.lastID = userID
	return m.response, m.err
}

type mockSSHKeyDeleteService struct {
	err    error
	lastID string
}

func (m *mockSSHKeyDeleteService) Delete(_ context.Context, sshKeyID string) error {
	m.lastID = sshKeyID
	return m.err
}

type mockOrganizationDeleteService struct {
	err      error
	lastName string
}

func (m *mockOrganizationDeleteService) Delete(_ context.Context, organization string) error {
	m.lastName = organization
	return m.err
}

type mockRunTaskDeleteReaderService struct {
	readResponse *tfe.RunTask
	readErr      error
	deleteErr    error
	lastReadID   string
	lastDeleteID string
}

func (m *mockRunTaskDeleteReaderService) Read(_ context.Context, runTaskID string) (*tfe.RunTask, error) {
	m.lastReadID = runTaskID
	return m.readResponse, m.readErr
}

func (m *mockRunTaskDeleteReaderService) Delete(_ context.Context, runTaskID string) error {
	m.lastDeleteID = runTaskID
	return m.deleteErr
}

type mockAgentPoolDeleteService struct {
	err    error
	lastID string
}

func (m *mockAgentPoolDeleteService) Delete(_ context.Context, agentPoolID string) error {
	m.lastID = agentPoolID
	return m.err
}

type mockNotificationDeleteService struct {
	err    error
	lastID string
}

func (m *mockNotificationDeleteService) Delete(_ context.Context, notificationConfigurationID string) error {
	m.lastID = notificationConfigurationID
	return m.err
}

type mockAuditTrailTokenDeleteService struct {
	err      error
	lastOrg  string
	lastOpts tfe.OrganizationTokenDeleteOptions
}

func (m *mockAuditTrailTokenDeleteService) DeleteWithOptions(_ context.Context, organization string, options tfe.OrganizationTokenDeleteOptions) error {
	m.lastOrg = organization
	m.lastOpts = options
	return m.err
}

type mockOrganizationUpdateService struct {
	response    *tfe.Organization
	err         error
	lastName    string
	lastOptions tfe.OrganizationUpdateOptions
}

func (m *mockOrganizationUpdateService) Update(_ context.Context, organization string, options tfe.OrganizationUpdateOptions) (*tfe.Organization, error) {
	m.lastName = organization
	m.lastOptions = options
	return m.response, m.err
}

type mockReservedTagKeyCreateService struct {
	response    *tfe.ReservedTagKey
	err         error
	lastOrg     string
	lastOptions tfe.ReservedTagKeyCreateOptions
}

func (m *mockReservedTagKeyCreateService) Create(_ context.Context, organization string, options tfe.ReservedTagKeyCreateOptions) (*tfe.ReservedTagKey, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}
