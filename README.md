# MEMPHIS PROTOCOL ADAPTER

[![Go](https://github.com/g41797/memphis-protocol-adapter/actions/workflows/go.yml/badge.svg)](https://github.com/g41797/memphis-protocol-adapter/actions/workflows/go.yml)

  This project is developing in accordance with [#849](https://github.com/memphisdev/memphis/issues/849)

  First developed adapter is **syslog-adapter**. 
  

syslog-adapter is based on 
- [syslogsidecar](https://github.com/g41797/syslogsidecar#readme)
- [sputnik](https://github.com/g41797/sputnik)

syslog-adapter consists of:
- syslog server - common part for all syslogsidecar based processes
- memphis specific plugins 

## Syslog server

 Supported RFCs:
  - [RFC3164](<https://tools.ietf.org/html/rfc3164>)
  - [RFC5424](<https://tools.ietf.org/html/rfc5424>)


  RFC3164 message consists of following symbolic parts:
  - priority
  - facility 
  - severity
  - timestamp
  - hostname
  - tag
  - **content**

  ### RFC5424

  RFC5424 message consists of following symbolic parts:
 - priority
 - facility 
 - severity
 - timestamp
 - hostname
 - version
 - app_name
 - proc_id
 - msg_id
 - structured_data
 - **message**

 ### Non-RFC parts

  syslogsidecar adds rfc of produced message:
  - Part name: "rfc"
  - Values: "RFC3164"|"RFC5424"

### Badly formatted messages

  syslogsidecar creates only one part for badly formatted message - former syslog message:
  - Part name: "data"
      
      
### Severities

  Valid severity levels and names are:
 - 0 emerg
 - 1 alert
 - 2 crit
 - 3 err
 - 4 warning
 - 5 notice
 - 6 info
 - 7 debug

  syslogsidecar filters messages by level according to value in configuration, e.g. for:
```json
{
  "SEVERITYLEVEL": 4,
  ...........
}
```
all messages with severity above 4 will be discarded. 


  ### Configuration

  Configuration of syslog server part of syslogsidecar is saved in the file syslogreceiver.json:
```json
{
    "SEVERITYLEVEL": 4,
    "ADDRTCP": "127.0.0.1:5141",
    "ADDRUDP": "127.0.0.1:5141",
    "UDSPATH": "",
    "ADDRTCPTLS": "127.0.0.1:5143",
    "CLIENT_CERT_PATH": "",
    "CLIENT_KEY_PATH ": "",
    "ROOT_CA_PATH": ""
}
```

### Links

- More complete description of [syslogsidecar](https://github.com/g41797/syslogsidecar#readme)


## Memphis Plugins

### Connector

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

Connector creates shared _*nats.Conn*_ for:
- periodic validation of connectivity with memphis
- loggers



### Producer

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
