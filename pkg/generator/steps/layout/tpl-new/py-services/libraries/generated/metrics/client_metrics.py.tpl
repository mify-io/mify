{{ .TplHeader }}

from prometheus_client import Counter, Histogram


class ClientMetrics:
    def __init__(self):
        self._request_duration = Histogram(
            namespace="service",
            subsystem="client_api",
            name="request_duration_seconds",
            documentation="Duration of request in ms",
            buckets=[0.02, 0.05, 0.1, 0.2, 0.5, 1, 5, 30, 60],
            labelnames=["target_service", "path"],
        )

        self._request_count = Counter(
            namespace="service",
            subsystem="client_api",
            name="requests_total",
            documentation="Total count of processed requests",
            labelnames=["target_service", "host", "path", "status"],
        )


    def report_request_end(self,
           status: int,
           duration: float,
           target_service: str,
           target_host: str,
           target_path: str):
        self._request_duration.labels(
                target_service,
                target_path).observe(duration)

        self._request_count.labels(
                target_service,
                target_host,
                target_path,
                str(status)).inc()
