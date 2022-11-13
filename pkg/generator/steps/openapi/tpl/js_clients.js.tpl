{{- .Header}}

{{- range .Clients}}
import {{.ClassName}} from '@/generated/api/clients/{{.ClientName}}'
{{- end}}

class Clients {
  constructor(ctx) {
    this.ctx = ctx
{{- range .Clients}}
    this._{{.PublicMethodName}} = new {{.ClassName}}.Api(
        new {{.ClassName}}.ApiClient(ctx.store.state.config),
    );
{{- end}}
  }

{{- range .Clients}}
  {{.PublicMethodName}}() {
    return this._{{.PublicMethodName}};
  }
{{- end}}

  static getConfigEnvMap() {
{{- range .Clients}}
      var {{.PublicMethodName}}Env = {{.ClassName}}.ApiClient.getConfigEnvName()
{{- end}}
      return {
{{- range .Clients}}
          [{{.PublicMethodName}}Env]: process.env[{{.PublicMethodName}}Env],
{{- end}}
      }

  }
}

export default Clients
