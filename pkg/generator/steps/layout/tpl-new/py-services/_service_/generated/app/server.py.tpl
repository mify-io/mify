{{.TplHeader}}

import asyncio

from aiohttp import web

from libraries.generated.metrics.server import MetricsServer

from {{.ServiceName}}.generated.openapi import app
from {{.ServiceName}}.generated.core.service_context import MifyServiceContext


class ServerConf:
    API_ENDPOINT: str = '{{.ApiEndpoint}}'
    MAINTENANCE_ENDPOINT: str = '{{.MaintenanceEndpoint}}'

    _env_mapping: dict = {
        'API_ENDPOINT': '{{.ApiEndpointEnv}}',
        'MAINTENANCE_ENDPOINT': '{{.MaintenanceEndpointEnv}}',
    }


class Server:
    def __init__(self, ctx: MifyServiceContext):
        self._config = ctx.static_config.get_config(ServerConf)
        self._context = ctx
        self._runners = []

    def _make_hostport(self, is_maintenance):
        conf = self._config.API_ENDPOINT
        conf_name = 'API_ENDPOINT'
        if is_maintenance:
            conf = self._config.MAINTENANCE_ENDPOINT
            conf_name = 'MAINTENANCE_ENDPOINT'
        hostport = conf.split(':')
        if len(hostport) < 2:
            raise ValueError(f"invalid {self._config._env_mapping[conf_name]} string, should be in format: [host]:port")
        return (hostport[0], hostport[1])

    async def _create_site(self, app, srv_name, host, port):
        self._context.logger.info(f"starting {srv_name} server", endpoint=f"{host}:{port}")
        runner = web.AppRunner(app)
        self._runners.append(runner)
        await runner.setup()
        site = web.TCPSite(runner, host, port)
        await site.start()
        return runner

    def _create_maintenance_app(self):
        metrics_srv = MetricsServer(self._context)

        app = web.Application()
        app.add_routes([
            web.get('/', metrics_srv.root_page),
            web.get('/metrics', metrics_srv.stats_page),
        ])
        return app

    async def _start_apps(self):
        api_app = app.create_server(self._context).app
        maintenance_app = self._create_maintenance_app()
        self._runners.append(await self._create_site(api_app, "api", *self._make_hostport(False)))
        self._runners.append(await self._create_site(maintenance_app, "maintenance", *self._make_hostport(True)))

    def _cleanup(self, loop):
        for runner in self._runners:
            loop.run_until_complete(runner.cleanup())

    def run(self):
        loop = asyncio.get_event_loop()
        loop.create_task(self._start_apps())
        try:
            loop.run_forever()
        except KeyboardInterrupt:
            pass
        finally:
           self._cleanup(loop)
