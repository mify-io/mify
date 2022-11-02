{{ .TplHeader }}

from libraries.generated.metrics.request_metrics import RequestMetrics

from .service_context import MifyServiceContext

_REQUEST_METRICS = RequestMetrics()

class MifyRequestContextBuilder:
    def __init__(self, service_ctx):
        self._service_ctx = service_ctx
        self._logger = service_ctx.logger
        self._metrics = _REQUEST_METRICS

    def set_protocol(self, protocol):
        self._protocol = protocol

    def set_request_id(self, request_id):
        self._request_id = request_id

    def set_url_path(self, url_path):
        self._url_path = url_path

    @property
    def logger(self):
        logger = self._logger.bind(
            request_id=self._request_id,
            proto=self._protocol,
            path=self._url_path,
        )
        return logger

    @property
    def service_ctx(self):
        return self._service_ctx

    @property
    def url_path(self):
        return self._url_path

    @property
    def metrics(self):
        return self._metrics

    def build(self, request):
        return MifyRequestContext(
            self._service_ctx,
            self._request_id,
            self.logger, # add
            request,
        )

class MifyRequestContext(MifyServiceContext):
    def __init__(self, service_ctx, request_id, logger, request):
        self._service_ctx = service_ctx
        self._request_id = request_id
        self._logger = logger
        self._request = request

    @property
    def service_context(self):
        return self._service_ctx

    @property
    def request_id(self):
        return self._request_id

    @property
    def logger(self):
        return self._logger

    @property
    def request(self):
        return self._request
