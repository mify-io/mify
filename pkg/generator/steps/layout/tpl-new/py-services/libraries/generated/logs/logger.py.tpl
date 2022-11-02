import structlog
import logging

def rename_event_key(_, __, ed):
    ed["msg"] = ed.pop("event")
    return ed

class MifyLoggerWrapper:
    def create_logger(service_context):
        structlog.configure(
            processors=[
                rename_event_key,
                # TODO: caller field unification with go service
                structlog.processors.CallsiteParameterAdder(
                    [
                        structlog.processors.CallsiteParameter.FILENAME,
                        structlog.processors.CallsiteParameter.LINENO,
                    ]
                ),
                structlog.processors.add_log_level,
                structlog.processors.StackInfoRenderer(),
                structlog.dev.set_exc_info,
                structlog.processors.TimeStamper(
                    fmt="ISO",
                    key="@timestamp",
                ),
                structlog.processors.JSONRenderer(),
            ],
            wrapper_class=structlog.make_filtering_bound_logger(logging.NOTSET),
            context_class=dict,
            logger_factory=structlog.PrintLoggerFactory(),
            cache_logger_on_first_use=False
        )
        logger = structlog.get_logger()
        logger = logger.bind(
            service_name=service_context.service_name,
            hostname=service_context.hostname,
            )
        return logger
