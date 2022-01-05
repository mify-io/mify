package mifyconfig

const (
	ServiceConfigName = "service.mify.yaml"
)

type ServiceLanguage string

const (
	ServiceLanguageUnknown ServiceLanguage = "unknown"
	ServiceLanguageGo      ServiceLanguage = "go"
	ServiceLanguageJs      ServiceLanguage = "js"
)

var LanguagesList = []ServiceLanguage{
	ServiceLanguageGo,
	ServiceLanguageJs,
}

type ServiceOpenAPIClientConfig struct {}

type ServiceOpenAPIConfig struct {
	Clients map[string]ServiceOpenAPIClientConfig `yaml:"clients,omitempty"`
}

type ServiceConfig struct {
	Language ServiceLanguage `yaml:"-"`

	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`

	OpenAPI ServiceOpenAPIConfig `yaml:"openapi,omitempty"`
}

type ServiceInfo struct {
	Name string
	ConfigPath string
}
