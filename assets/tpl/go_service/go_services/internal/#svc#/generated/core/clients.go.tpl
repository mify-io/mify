{{- .Workspace.TplHeader}}

package core

import (
	{{- range .OpenAPI.Clients}}
	"{{$.GoModule}}/internal/{{$.ServiceName}}/generated/api/clients/{{.ClientName}}"
	{{- end}}

	{{if .OpenAPI.Clients }}"{{$.GoModule}}/internal/pkg/generated/metrics"{{end}}
)

type MifyServiceClients struct {
	{{- range .OpenAPI.Clients}}
	{{.PrivateFieldName}} *{{.PackageName}}.APIClient
	{{- end}}
}

func NewMifyServiceClients(ctx *MifyServiceContext) (*MifyServiceClients, error) {
	{{- range .OpenAPI.Clients}}
	{{.PrivateFieldName}} := {{.PackageName}}.NewAPIClient(metrics.NewClientMetrics(), {{.PackageName}}.NewConfiguration(ctx.StaticConfig()))
	{{- end}}

	clients := &MifyServiceClients {
		{{- range .OpenAPI.Clients}}
		{{.PrivateFieldName}}: {{.PrivateFieldName}},
		{{- end}}
	}

	return clients, nil
}

{{- range .OpenAPI.Clients}}
func (c *MifyServiceClients) {{.PublicMethodName}}() *{{.PackageName}}.APIClient {
	return c.{{.PrivateFieldName}}
}
{{- end}}
