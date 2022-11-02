{{ .TplHeader }}
from dataclasses import dataclass

from prometheus_client import Counter, Summary, Histogram


@dataclass
class RequestInfo:
    service_name: str
    hostname: str
    url_path: str


class RequestMetrics:
    def __init__(self):
        self._request_size = Summary(
            namespace="service",
            subsystem="api",
            name="request_size_bytes",
            documentation="Size of input request",
            #  objectives={0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
            labelnames=["service", "host_name", "path"],)

        self._response_size = Summary(
            namespace="service",
            subsystem="api",
            name="response_size_bytes",
            documentation="Size of response",
            #  objectives={0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
            labelnames=["service", "host_name", "path"],
        )

        self._request_duration = Histogram(
            namespace="service",
            subsystem="api",
            name="request_duration_seconds",
            documentation="Duration of request in ms",
            buckets=[0.02, 0.05, 0.1, 0.2, 0.5, 1, 5, 30, 60],
            labelnames=["service", "host_name", "path"],
        )

        self._request_count = Counter(
            namespace="service",
            subsystem="api",
            name="requests_total",
            documentation="Total count of processed requests",
            labelnames=["service", "host_name", "path", "status"],
        )

    def report_request_end(
       self, req_info: RequestInfo, status: int,
       duration: float, requestSizeBytes: int, responseSizeBytes: int):
        self._request_size.labels(
            req_info.service_name,
            req_info.hostname,
            req_info.url_path).observe(float(requestSizeBytes))

        self._response_size.labels(
            req_info.service_name,
            req_info.hostname,
            req_info.url_path).observe(float(responseSizeBytes))

        self._request_duration.labels(
            req_info.service_name,
            req_info.hostname,
            req_info.url_path).observe(duration)

        self._request_count.labels(
            req_info.service_name,
            req_info.hostname,
            req_info.url_path,
            str(status)).inc()
