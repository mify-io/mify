{{- .Header}}

{{- range .Clients}}
import {{.PublicMethodName}} from '@/generated/api/clients/{{.ClientName}}'
{{- end}}

export default {
{{- range .Clients}}
	{{.PublicMethodName}},
{{- end}}
}
