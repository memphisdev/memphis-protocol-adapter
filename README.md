# MEMPHIS PROTOCOL ADAPTER

[![Go](https://github.com/g41797/memphis-protocol-adapter/actions/workflows/go.yml/badge.svg)](https://github.com/g41797/memphis-protocol-adapter/actions/workflows/go.yml)

  This project is developing in accordance with [#849](https://github.com/memphisdev/memphis/issues/849)

  First developed adapter is **syslog-adapter**. 
  
  It is based on [syslogsidecar framework](https://github.com/g41797/syslogsidecar#readme).

  Implementation for memphis consists of 3 plugins:
  - connector
  - producer
  - consumer (used for the tests)


## Connector

Configuration file: connector.json
```json
{
    "MEMPHIS_ADDR": "localhost:6666",
    "MEMPHIS_CLIENT":"MEMPHIS HTTP LOGGER",
    "USER_PASS_BASED_AUTH": true,
    "ROOT_USER": "root",
    "ROOT_PASSWORD": "memphis",
    "CONNECTION_TOKEN": "memphis",
    "CLIENT_CERT_PATH": "",
    "CLIENT_KEY_PATH ": "",
    "ROOT_CA_PATH": "",
    "CLOUD_ENV": false,
    "DEBUG": true,
    "DEV_ENV": true
}
```

Part of configuration is placed within docker-compose.yml:
```yml
    environment:
          - MEMPHIS_ADDR=memphis:6666
```

Connector creates sharable _*nats.Conn*_ for:
- periodic validation of connectivity with memphis
- loggers



## Producer

Configuration file: syslogproducer.json
```json
{
    "MEMPHIS_HOST": "127.0.0.1",
    "MEMPHIS_USER": "root",
    "MEMPHIS_PSWRD": "memphis",
    "PRODUCER": "syslog-adapter",
    "STATION": "syslog",
    "RETENTIONTYPE":"MaxMessageAgeSeconds",
    "RETENTIONVALUE":600
}
```

Part of configuration is placed within docker-compose.yml:
```yml
    environment:
          - MEMPHIS_HOST=memphis
```

Producer uses dedicated _*memphis.Conn*_.

syslog messages are produced to memphis as *MsgHeaders* with empty payload:
```go
err := mpr.producer.Produce("", memphis.MsgHeaders(hdrs))
```



