{{- .Header}}

package core

import (
	{{- range .Clients}}
	"{{.IncludePath}}"
	{{- end}}

	{{if .MetricsIncludePath }}"{{.MetricsIncludePath}}"{{end}}
)

type MifyServiceClients struct {
	{{- range .Clients}}
	{{.PrivateFieldName}} *{{.PackageName}}.APIClient
	{{- end}}
}

func NewMifyServiceClients(ctx *MifyServiceContext) (*MifyServiceClients, error) {
	{{- range .Clients}}
	{{.PrivateFieldName}} := {{.PackageName}}.NewAPIClient(metrics.NewClientMetrics(), {{.PackageName}}.NewConfiguration(ctx.StaticConfig()))
	{{- end}}

	clients := &MifyServiceClients {
		{{- range .Clients}}
		{{.PrivateFieldName}}: {{.PrivateFieldName}},
		{{- end}}
	}

	return clients, nil
}

{{- range .Clients}}
func (c *MifyServiceClients) {{.PublicMethodName}}() *{{.PackageName}}.APIClient {
	return c.{{.PrivateFieldName}}
}
{{- end}}
