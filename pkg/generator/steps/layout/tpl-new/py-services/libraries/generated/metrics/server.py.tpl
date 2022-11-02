from typing import Callable

from aiohttp import web

from prometheus_client import CONTENT_TYPE_LATEST, REGISTRY, generate_latest
from prometheus_client.openmetrics import exposition as openmetrics


_ROOT_CONTENT = '<html><body><a href="/metrics">Metrics</a></body></html>'


def _choose_generator(accept_header: Optional[str]) -> tuple[Callable, str]:
    accept_header = accept_header or ""
    for accepted in accept_header.split(","):
        if accepted.split(";")[0].strip() == "application/openmetrics-text":
            return (
                openmetrics.generate_latest,
                openmetrics.CONTENT_TYPE_LATEST,
            )

    return generate_latest, CONTENT_TYPE_LATEST


class MetricsServer:
    def __init__(self, context):
        self._context = context

    async def stats_page(self, request: web.Request) -> web.Response:
        generate, content_type = _choose_generator(request.headers.get("Accept"))

        rsp = web.Response(body=generate(REGISTRY))
        # This is set separately because aiohttp complains about `;` in
        # content_type thinking it means there's also a charset.
        # cf. https://github.com/aio-libs/aiohttp/issues/2197
        rsp.content_type = content_type

        return rsp

    async def root_page(self, request: web.Request) -> web.Response:
        return web.Response(text=_ROOT_CONTENT, content_type="text/html")
