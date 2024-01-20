{{.TplHeader}}
# vim: set ft=python:

from {{.Workspace.MifyGeneratedServicePackage}}.core.service_context import MifyServiceContext
from {{.ServiceName}}.app.service_extra import ServiceExtra

from .server import Server

class MifyServiceApp:
    def __init__(self):
        self._service_context = MifyServiceContext('{{.ServiceName}}').with_extra(ServiceExtra)
        self._server = Server(self._service_context)

    def run(self):
        self._server.run()
        self._service_context.close()
