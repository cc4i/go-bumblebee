# Air

Air service is backend for Gate service to retrieve air quality data, which will fire request to [AQI](https://api.waqi.info) and get data back.

## Features

- /ping
- /air/city/:city

## Environments

AQI_SERVER_URL=https://127.0.0.1
AQI_SERVER_TOKEN=
IP_STACK_SERVER_URL=http://127.0.0.1
IP_STACK_SERVER_TOKEN=

REDIS_SERVER_ADDRESS=127.0.0.1:6379
REDIS_SERVER_PASSWORD=

AIR_VERSION=v0.1.0