# MEMPHIS PROTOCOL ADAPTER

  This project is developing in accordance with [#849](https://github.com/memphisdev/memphis/issues/849)

  First developed adapter is **syslog-adapter**

## syslog-adapter

  syslog-adapter:
  - receives messages intended for [syslogd](https://linux.die.net/man/8/syslogd)
  - sends them (produces using **MsgHeaders**) to the broker

  Supported RFCs:
  - [RFC3164](<https://tools.ietf.org/html/rfc3164>)
  - [RFC5424](<https://tools.ietf.org/html/rfc5424>)

  User friendly description of syslogformat:[Analyze syslog messages](https://blog.datalust.co/seq-input-syslog/)


  ### RFC3164

  RFC3164 is oldest syslog RFC, syslog-adapter supports it for old syslogd clients.

  RFC3164 message consists of following parts:
  - priority
	- facility 
	- severity
	- timestamp
	- hostname
  - tag
	- content

  Example: after send of RFC3164 log message with current time as *content* to syslog-adapter
```bash
  timestamp=$(date +%d-%m-%Y_%H-%M-%S)
  logger  --rfc3164 --server 127.0.0.1 --port 5141 --priority user.alert  $timestamp
```

  it will be represented in JSON format on broker side as:
```json
{
  "MsgHeaders": {
    "client": [
      ""
    ],
    "content": [
      "23-08-2023_08-53-54"
    ],
    "facility": [
      "1"
    ],
    "hostname": [
      "BLKF"
    ],
    "priority": [
      "9"
    ],
    "rfc": [
      "RFC3164"
    ],
    "severity": [
      "1"
    ],
    "tag": [
      "g41797"
    ],
    "timestamp": [
      "2023-08-23 08:53:54 +0000 UTC"
    ],
    "tls_peer": [
      ""
    ]
  }
}
```

  ### RFC5424

  RFC5424 message consists of following parts:
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
	- message

  Example: after send of RFC5424 log message with current time as *message* to syslog-adapter
```bash
  timestamp=$(date +%d-%m-%Y_%H-%M-%S)
  logger  --rfc5424 --server 127.0.0.1 --port 5141 --priority user.alert  $timestamp
```

  it will be represented in JSON format  on broker side as:
```json
{
  "MsgHeaders": {
    "app_name": [
      "g41797"
    ],
    "client": [
      ""
    ],
    "facility": [
      "1"
    ],
    "hostname": [
      "BLKF"
    ],
    "message": [
      "23-08-2023_08-53-54"
    ],
    "msg_id": [
      "-"
    ],
    "priority": [
      "9"
    ],
    "proc_id": [
      "-"
    ],
    "rfc": [
      "RFC5424"
    ],
    "severity": [
      "1"
    ],
    "structured_data": [
      "[timeQuality tzKnown=\"1\" isSynced=\"1\" syncAccuracy=\"45460\"]"
    ],
    "timestamp": [
      "2023-08-23 05:53:54.040825 +0000 UTC"
    ],
    "tls_peer": [
      ""
    ],
    "version": [
      "1"
    ]
  }
}
```
  ### Non-RFC parts

  syslog-adapter adds rfc of produced message to standard parts:

```json
{
  "MsgHeaders": {
  ------------
    "rfc": [
      "RFC3164"
    ],
  ------------
}
}

{
  "MsgHeaders": {
  ------------
    "rfc": [
      "RFC5424"
    ],
  ------------
}
}

```

Looks that *client* and *tls_peer* are added by memphis infrastructure.


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

```bash
 <process> -cf <path to folder with configuration files>
```


Example of running in vscode terminal
```bash
 ./syslog-adapter -cf ./cmd/syslog-adapter/conf/
```

