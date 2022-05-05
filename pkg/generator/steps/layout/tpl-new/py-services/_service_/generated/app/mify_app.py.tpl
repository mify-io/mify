{{.TplHeader}}

from {{.ServiceName}}.generated.core.service_context import MifyServiceContext

from .server import Server

class MifyServiceApp:
    def __init__(self):
        self._service_context = MifyServiceContext('{{.ServiceName}}')
        self._server = Server(self._service_context)

    def run(self):
        self._server.run()
        self._service_context.close()
