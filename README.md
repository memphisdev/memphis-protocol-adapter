# MEMPHIS PROTOCOL ADAPTER

## Structure of the repository

```bash
memphis-protocol-adapter
.
├── cmd
│   ├── protocol-adapter
│   │   └── main.go
│   ├── rest-gateway
│   │   └── main.go
│   └── syslog-adapter
│       ├── conf
│       │   ├── blocks.json
│       │   ├── connector.json
│       │   ├── syslogproducer.json
│       │   └── syslogreceiver.json
│       └── main.go
├── go.mod
├── go.sum
├── LICENSE
├── pkg
│   ├── adapter
│   │   ├── config.go
│   │   ├── connector.go
│   │   ├── logger.go
│   │   └── runner.go
│   ├── rest
│   └── syslogblocks
│       ├── msgproducer.go
│       ├── producer.go
│       ├── receiver.go
│       └── server.go
└── README.md
```

### `/cmd`

Applications for this project:
* syslog-adapter 
* rest-gateway
* protocol-adapter


#### `/cmd/syslog-adapter`

##### `/cmd/syslog-adapter/conf`


#### `/cmd/rest-gateway`

  Placeholder for the rest gateway 

#### `/cmd/protocol-adapter`

  Placeholder for the process uniting all adapters


### `/pkg`

Packages of the project
* rest - placeholder for the rest gateway 
* syslog - code for the syslog adapter(based on [Memphis Go SDK](https://github.com/memphisdev/memphis.go) and [sputnik](https://github.com/g41797/sputnik))
* adapter - runtime 

