# MEMPHIS PROTOCOL ADAPTER

## Structure of the repository

```bash
memphis-protocol-adapter
├── cmd
│   ├── protocol-adapter
│   │   ├── blocks.yaml
│   │   ├── config.json
│   │   └── main.go
│   ├── rest-gateway
│   │   ├── blocks.yaml
│   │   ├── config.json
│   │   └── main.go
│   └── syslog-adapter
│       ├── blocks.yaml
│       ├── config.json
│       └── main.go
├── go.mod
├── LICENSE
├── Makefile
├── pkg
│   ├── rest
│   │   ├── producer.go
│   │   └── receiver.go
│   ├── server
│   │   ├── config.go
│   │   ├── connector.go
│   │   └── logger.go
│   └── syslog
│       ├── producer.go
│       └── receiver.go
├── README.md
└── scripts
```

### `/cmd`

Applications for this project:
* syslog-adapter 
* rest-gateway
* protocol-adapter


#### `/cmd/syslog-adapter`

Executable includes support for *syslog* (see blocks.yaml)
```yaml
Blocks:
- Name: syslogreceiver
  Responsibility: syslogreceiver
- Name: syslogproducer
  Responsibility: syslogproducer
```

It includes 2 **Blocks**:
* syslogreceiver - receives syslogs from the clients
  * implemented in pkg/syslog/receiver.go
* syslogproducer - acts as memphis client, produces syslog events
   * implemented in pkg/syslog/producer.go

#### `/cmd/rest-gateway`

Executable includes support for *REST-HTTP* (see blocks.yaml)
```yaml
Blocks:
- Name: restreceiver
  Responsibility: restreceiver
- Name: restproducer
  Responsibility: restproducer
```

It includes 2 **Blocks**:
* restreceiver - receives HTTP REST requests from the clients
  * implemented in pkg/rest/receiver.go
* restproducer - acts as memphis client
   * implemented in pkg/rest/producer.go


#### `/cmd/protocol-adapter`

Executable includes support for *syslog* and *rest* (see blocks.yaml)
```yaml
Blocks:
- Name: syslogreceiver
  Responsibility: syslogreceiver
- Name: syslogproducer
  Responsibility: syslogproducer
- Name: restreceiver
  Responsibility: restreceiver
- Name: restproducer
  Responsibility: restproducer  
```

### `/pkg`

Packages of the project
* rest - code for the rest gateway (based on [sputnik](https://github.com/g41797/sputnik))
* syslog - code for the syslog adapter(based on [sputnik](https://github.com/g41797/sputnik))
* server - code for working with memphis server(based on [Memphis Go SDK](https://github.com/g41797/memphis.go))

