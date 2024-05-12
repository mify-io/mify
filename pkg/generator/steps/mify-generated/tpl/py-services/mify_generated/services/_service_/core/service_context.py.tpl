{{.TplHeader}}
# vim: set ft=python:

import structlog
import logging
import socket
import os

from {{.Workspace.MifyGeneratedCommonPackage}}.logs.logger import MifyLoggerWrapper
from {{.Workspace.MifyGeneratedCommonPackage}}.configs.static import MifyStaticConfig
from {{.Workspace.MifyGeneratedCommonPackage}}.configs.dynamic import MifyDynamicConfig

from .clients import MifyServiceClients

class MifyServiceContext:
    def __init__(self, service_name):
        self._hostname = socket.gethostname()
        self._service_name = service_name
        self._logger = MifyLoggerWrapper.create_logger(self)
        self._static_config = MifyStaticConfig(os.environ)
        self._dynamic_config = MifyDynamicConfig(self._static_config)
        self._clients = MifyServiceClients(self)

    def close(self):
        self._clients.close()
        self._dynamic_config.close()

    def with_extra(self, extra):
        self._extra = extra(self)
        return self

    @property
    def hostname(self):
        return self._hostname

    @property
    def service_name(self):
        return self._service_name

    @property
    def logger(self):
        return self._logger

    @property
    def static_config(self):
        return self._static_config

    @property
    def dynamic_config(self):
        return self._dynamic_config

    @property
    def service_extra(self):
        return self._extra

    @property
    def clients(self):
        return self._clients
