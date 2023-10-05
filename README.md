# MEMPHIS PROTOCOL ADAPTER

[![Go](https://github.com/g41797/memphis-protocol-adapter/actions/workflows/go.yml/badge.svg)](https://github.com/g41797/memphis-protocol-adapter/actions/workflows/go.yml)

  This project is developing in accordance with [#849](https://github.com/memphisdev/memphis/issues/849)

  First developed adapter is **syslog-adapter**

## syslog-adapter

syslog-adapter is based on [syslogsidecar framework](https://github.com/g41797/syslogsidecar#readme)

### Command line

Example of running in vscode terminal
```bash
 ./syslog-adapter -cf ./cmd/syslog-adapter/conf/
```

### e2e tests

Functionality: asynchronously
- send 1000000 syslog messages via one TCP/IP connection to syslogsidecar
- receive messages and forward to the broker
- consume messages 
- compare 
- print report

Build and run under vscode:
```bash
go clean -cache -testcache
go build ./cmd/syslog-e2e/
./syslog-e2e -cf ./cmd/syslog-e2e/conf/
```
Required memphis services will be started automatically according to *conf/docker-compose.yml* file

