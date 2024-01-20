{{.TplHeader}}
# vim: set ft=python:
from {{.Workspace.MifyGeneratedCommonPackage}}.metrics.client_metrics import ClientMetrics

{{- range .Model.Clients}}
from {{.ImportPath}} import ApiClient as {{.ApiClientName}}, Configuration as {{.ConfigurationName}}
{{- end}}

class MifyServiceClients:
    def __init__(self, service_context):
        metrics = ClientMetrics()
        {{- range .Model.Clients}}
        self._{{.PropertyName}} = {{.ApiClientName}}(metrics, {{.ConfigurationName}}(service_context.static_config))
        {{- end}}
        pass

    def close(self):
        {{- range .Model.Clients}}
        self._{{.PropertyName}}.close()
        {{- end}}
        pass

    {{- range .Model.Clients}}
    @property
    def {{.PropertyName}}(self):
        return self._{{.PropertyName}}
    {{- end}}
