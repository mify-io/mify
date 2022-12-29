{{.TplHeader}}
# vim: set ft=python:
from libraries.generated.metrics.client_metrics import ClientMetrics

{{- range .Clients}}
from {{.ImportPath}} import ApiClient as {{.ApiClientName}}, Configuration as {{.ConfigurationName}}
{{- end}}

class MifyServiceClients:
    def __init__(self, service_context):
        metrics = ClientMetrics()
        {{- range .Clients}}
        self._{{.PropertyName}} = {{.ApiClientName}}(metrics, {{.ConfigurationName}}(service_context.static_config))
        {{- end}}
        pass

    def close(self):
        {{- range .Clients}}
        self._{{.PropertyName}}.close()
        {{- end}}
        pass

    {{- range .Clients}}
    @property
    def {{.PropertyName}}(self):
        return self._{{.PropertyName}}
    {{- end}}
