connexion[aiohttp,swagger-ui] >= 2.6.0; python_version>="3.6"
# 2.3 is the last version that supports python 3.5
connexion[aiohttp,swagger-ui] <= 2.3.0; python_version=="3.5" or python_version=="3.4"
# connexion requires werkzeug but connexion < 2.4.0 does not install werkzeug
# we must peg werkzeug versions below to fix connexion
# https://github.com/zalando/connexion/pull/1044
werkzeug == 0.16.1; python_version=="3.5" or python_version=="3.4"
swagger-ui-bundle == 0.0.6
aiohttp_jinja2 == 1.5.0
aiohttp_cors >= 0.7.0
structlog >= 22.1.0
prometheus-client >= 0.14.1
certifi >= 14.05.14
frozendict >= 2.0.3
python_dateutil >= 2.5.3
setuptools >= 21.0.0
urllib3 >= 1.15.1
typing-extensions >= 4.4.0
python-consul2 >= 0.1.5
