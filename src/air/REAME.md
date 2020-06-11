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
    
    MYSQL_SERVER_ADDRESS=127.0.0.1:6379
    MYSQL_SERVER_PASSWORD=
    
    AIR_VERSION=v0.1.0
    
    OVERRIDE_VERSION=v1 # for testing purpose in k8s

## Redis

## Table

    CREATE TABLE `aqi_call_logs` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `method` VARCHAR(64) NOT NULL,
        `url` VARCHAR(256) NOT NULL,
        `created` DATE NOT NULL,
        `request` VARCHAR(256),
        `response_code` INT(10),
        `response` VARCHAR(256),
        PRIMARY KEY (`uid`)
    );
    
    CREATE TABLE `air` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `namesapce` VARCHAR(64) NOT NULL,
        `pod` VARCHAR(64) NOT NULL,
        `created` DATE NOT NULL,
        PRIMARY KEY (`uid`)
    );    
    
 