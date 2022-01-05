package mifyconfig

const (
	WorkspaceConfigName = "workspace.mify.yaml"

	GoServicesRoot      = "go_services"
	JsServicesRoot      = "js_services"
)

type WorkspaceConfig struct {
	WorkspaceName string `yaml:"workspace_name"`
	GitHost       string `yaml:"git_host"`
	GitNamespace  string `yaml:"git_namespace"`
	GitRepository string `yaml:"git_repository,omitempty"`
}
