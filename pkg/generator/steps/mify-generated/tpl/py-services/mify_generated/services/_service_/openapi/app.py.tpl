{{.TplHeader}}
# vim: set ft=python:

import os
import connexion
import aiohttp_cors
import uuid
import time
import traceback

from aiohttp import web
from {{.Workspace.MifyGeneratedCommonPackage}}.metrics.request_metrics import RequestInfo
from {{.Workspace.MifyGeneratedServicePackage}}.core.request_context import MifyRequestContextBuilder


async def request_context(request: web.Request):
    service_ctx = request.config_dict['service_context']
    request_ctx_builder = MifyRequestContextBuilder(service_ctx)
    request_ctx_builder.set_request_id(str(uuid.uuid4()))
    request_ctx_builder.set_url_path(request.path)
    request_ctx_builder.set_protocol(request.url.scheme)
    request.app['request_context_builder'] = request_ctx_builder
    return request_ctx_builder


def trim_data(data, limit):
    if data is None:
        return '(null)', 0
    size = len(data)
    if size > limit:
        data = data[:limit] + '...'
    return data, size


async def log_request(ctx: MifyRequestContextBuilder, request: web.Request):
    body, orig_size = trim_data(await request.text(), limit=1024)
    ctx.logger.bind(request_body=body, size=orig_size).info("started processing request")
    return time.time()


async def log_response(ctx: MifyRequestContextBuilder, response: web.Response, start_time):
    elapsed = time.time() - start_time
    body, orig_size = trim_data(response.body, limit=1024)
    ctx.logger.bind(
            status=response.status,
            response_body=body,
            size=orig_size,
            elapsed_sec=elapsed).info("finished processing request")


async def write_metrics(
        ctx: MifyRequestContextBuilder, request: web.Request, response: web.Response, start_time):
    elapsed = time.time() - start_time
    ctx.metrics.report_request_end(
        RequestInfo(
            service_name=ctx.service_ctx.service_name,
            hostname=ctx.service_ctx.hostname,
            url_path=ctx.url_path,
        ), response.status, elapsed, len(await request.text()), response.body_length)


@web.middleware
async def middleware_handler(
        request: web.Request, handler: any) -> web.Response:
    request_ctx_builder = await request_context(request)
    start_time = await log_request(request_ctx_builder, request)

    response = None
    try:
        response = await handler(request)
    except Exception as e:
        trace = traceback.format_exc()
        request_ctx_builder.logger.bind(trace=trace).info("unhandled exception in request")
        response = web.json_response({
            "status": 500,
            "message": str(e),
        }, status=500)

    await log_response(request_ctx_builder, response, start_time)
    await write_metrics(request_ctx_builder, request, response, start_time)
    return response
    #  headers = request.headers
    #  x_auth_token = headers.get("X-Token")
    #  app_id = headers.get("X-AppId")
    #  user, success = security.Security().authorize(x_auth_token)
    #  if not success:
        #  return web.json_response(status=401, data={
            #  "error": {
                #  "message": ("Not authorized. Reason: {}"
                            #  )
            #  }
        #  })

def create_server(service_context):
    options = {
        "swagger_ui": True,
        'middlewares': [middleware_handler],
        }
    specification_dir = os.path.join(os.path.dirname(__file__), 'openapi')
    app = connexion.AioHttpApp(__name__,
            specification_dir=specification_dir,
            options=options,
            only_one_api=True,
    )

    app.add_api('openapi.yaml',
                arguments={'title': '{{.ServiceName}}'},
                pythonic_params=True,
                pass_context_arg_name='request')
    app.app['service_context'] = service_context

    # Enable CORS for all origins.
    cors = aiohttp_cors.setup(app.app, defaults={
        "*": aiohttp_cors.ResourceOptions(
            allow_credentials=True,
            expose_headers="*",
            allow_headers="*",
        )
    })

    # Register all routers for CORS.
    for route in list(app.app.router.routes()):
        cors.add(route)
    return app
