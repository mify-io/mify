{{- .Header}}

{{- range .Clients}}
import {{.ClassName}} from '../api/clients/{{.ClientName}}'
{{- end}}

class Clients {
  constructor(config) {
{{- range .Clients}}
    this._{{.PublicMethodName}} = new {{.ClassName}}.Api(
        new {{.ClassName}}.ApiClient(config),
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
