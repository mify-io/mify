{{- .Workspace.TplHeader}}

{{- range .OpenAPI.Clients}}
import {{.PublicMethodName}} from '@/generated/api/clients/{{.ClientName}}'
{{- end}}

export default {
{{- range .OpenAPI.Clients}}
	{{.PublicMethodName}},
{{- end}}
}
