# MEMPHIS PROTOCOL ADAPTER

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
- send 1000000 syslog messages
- consume messages 
- compare 
- print report

Build under vscode:
```bash
go clean -cache -testcache
go build ./rookie2e/cmd/syslog-re2e/
```
Run tests:
- Memphis DB
- Memphis Broker
- syslog-adapter
- syslog-re2e:
```bash
# under vscode terminal
./syslog-re2e -cf ./rookie2e/cmd/syslog-re2e/conf/
```


