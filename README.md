# MEMPHIS PROTOCOL ADAPTER

  This project is developing in accordance with [#849](https://github.com/memphisdev/memphis/issues/849)

  First developed adapter is **syslog-adapter**

## syslog-adapter

  syslog-adapter:
  - receives logs intended for [syslogd](https://linux.die.net/man/8/syslogd)
  - produces messages to the broker as **MsgHeaders**
     
  Supported RFCs:
  - [RFC3164](<https://tools.ietf.org/html/rfc3164>)
  - [RFC5424](<https://tools.ietf.org/html/rfc5424>)

  User friendly description of syslogformat:[Analyze syslog messages](https://blog.datalust.co/seq-input-syslog/)


  ### RFC3164

  RFC3164 is oldest syslog RFC, syslog-adapter supports it for old syslogd clients.

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

  syslog-adapter adds rfc of produced message to standard parts:
  - Part name: "rfc"
  - Values: "RFC3164"|"RFC5424"

  ### Facilities and severities

  Valid facility names are:
  - auth
  - authpriv for security information of a sensitive nature
  - cron
  - daemon
  - ftp
  - kern
  - lpr
  - mail
  - news
  - syslog
  - user
  - uucp
  - local0-local7

    Valid severity levels and names are:

 - 0 emerg
 - 1 alert
 - 2 crit
 - 3 err
 - 4 warning
 - 5 notice
 - 6 info
 - 7 debug

  syslog-adapters filters messages by level according to value in configuration:
```json
{
  "SEVERITYLEVEL": 4,
  ...........
}
```
All messages with severity above 4 will be discarded. 


  ### Configuration

  Configuration of receiver part of syslog-adapter is saved in the file syslogreceiver.json:
```json
{
    "SEVERITYLEVEL": 4,
    "ADDRTCP": "127.0.0.1:5141",
    "ADDRUDP": "0.0.0.0:5141",
    "UDSPATH": "",
    "ADDRTCPTLS": "127.0.0.1:5143",
    "CLIENT_CERT_PATH": "",
    "CLIENT_KEY_PATH ": "",
    "ROOT_CA_PATH": ""
}
```
and related go struct:
```go
type SyslogConfiguration struct {
	// The Syslog Severity level ranges between 0 to 7.
	// Each number points to the relevance of the action reported.
	// From a debugging message (7) to a completely unusable system (0):
	//
	//	0		Emergency: system is unusable
	//	1		Alert: action must be taken immediately
	//	2		Critical: critical conditions
	//	3		Error: error conditions
	//	4		Warning: warning conditions
	//	5		Notice: normal but significant condition
	//	6		Informational: informational messages
	//	7		Debug: debug-level messages
	//
	// Log with severity above value from configuration will be discarded
	// Examples:
	// -1 - all logs will be discarded
	// 5  - logs with severities 6(Informational) and 7(Debug) will be discarded
	// 7  - all logs will be processed
	SEVERITYLEVEL int

	// IPv4 address of TCP listener.
	// For empty string - don't use TCP
	// e.g "0.0.0.0:5141" - listen on all adapters, port 5141
	// "127.0.0.1:5141" - listen on loopback "adapter"
	ADDRTCP string

	// IPv4 address of UDP receiver.
	// For empty string - don't use UDP
	// Usually "0.0.0.0:5141" - receive from all adapters, port 5141
	// "127.0.0.1:5141" - receive from loopback "adapter"
	ADDRUDP string

	// Unix domain socket name - actually file path.
	// For empty string - don't use UDS
	// Regarding limitations see https://man7.org/linux/man-pages/man7/unix.7.html
	UDSPATH string

	// TLS section: Listening on non empty ADDRTCPTLS will start only
	// for valid tls configuration (created using last 3 parameters)
	ADDRTCPTLS       string
	CLIENT_CERT_PATH string
	CLIENT_KEY_PATH  string
	ROOT_CA_PATH     string
}
```

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


