package syslogblocks

import (
	syslogd "gopkg.in/mcuadros/go-syslog.v2"
)

type SyslogConfiguration struct {
	// TCP listener port. For 0 - don't use TCP
	PORTTCP int

	// Add after solving tls cert. probleb PORTTCPTLS int

	// UDP receiver port, For 0 - don't use UDP
	PORTUDP int

	// Unix domain socket name - actually file path.
	// For empty string - don't use UDS
	// Regarding limitations see https://man7.org/linux/man-pages/man7/unix.7.html
	UDSPATH string

	// The Syslog Severity level ranges between 0 to 7.
	// Each number points to the relevance of the action reported.
	// From a debugging message (7) to a completely unusable system (0)
	// Log with severity above value from configuration will be discarded
	// Examples:
	// -1 - all logs will be discarded
	// 5  - logs with severities 6(informational) and 7(debugging) will be discarded
	// 7  - all logs will be processed
	SEVERITYLEVEL int
}

type server struct {
	config  SyslogConfiguration
	syslogd *syslogd.Server
}

func (s *server) init(conf SyslogConfiguration) {
	s.config = conf
	s.syslogd = syslogd.NewServer()
}
