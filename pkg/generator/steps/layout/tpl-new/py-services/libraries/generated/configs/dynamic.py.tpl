# vim: set ft=python:
import asyncio
import consul
import consul.aio
import json

class MifyDynamicConfig:
    _CONFIG_PATH = 'config/'

    class ConsulConf:
        CONSUL_PORT: int = 8500
        CONSUL_HOST: str = '127.0.0.1'

    def __init__(self, static_config):
        config = static_config.get_config(self.ConsulConf)

        self._configs = {}
        self._loop = asyncio.get_event_loop()
        self._consul = consul.aio.Consul(
                host=config.CONSUL_HOST,
                port=config.CONSUL_PORT,
                loop=self._loop)
        self._consul_sync = consul.Consul(
                host=config.CONSUL_HOST,
                port=config.CONSUL_PORT)
        self._pollers = []
        pass

    async def _poll_config(self, conf_class: type):
        conf_name = conf_class.__name__
        index = None
        while True:
            index, data = await self._consul.kv.get(self._CONFIG_PATH + conf_name, index=index)
            self._configs[conf_name] = self._update_config(conf_class, data)

    def _update_config(self, conf_class, data):
        if data is None:
            return conf_class()
        data = json.loads(data['Value'])
        return conf_class(**data)

    def _add_or_get_config(self, conf_class: type):
        conf_name = conf_class.__name__
        if conf_name in self._configs:
            return self._configs[conf_name]

        _, data = self._consul_sync.kv.get(self._CONFIG_PATH + conf_name)
        self._configs[conf_name] = self._update_config(conf_class, data)

        self._pollers.append(self._loop.create_task(self._poll_config(conf_class)))
        return self._configs[conf_name]

    def close(self):
        for poller in self._pollers:
            poller.cancel()

    def get_config(self, conf_class: type):
        return self._add_or_get_config(conf_class)
