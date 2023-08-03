package syslogblocks_test

import (
	"testing"

	"github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
)

func Test_Init_Finish(t *testing.T) {

	srv := syslogblocks.Server{}
	defer stopServer(t, &srv)
	err := srv.Init(defaultServerConfiguration())
	if err != nil {
		t.Errorf("Init error %v", err)
	}
}

func stopServer(t *testing.T, srv *syslogblocks.Server) {
	err := srv.Finish()
	if err != nil {
		t.Errorf("stop server error %v", err)
	}
}

func defaultServerConfiguration() syslogblocks.SyslogConfiguration {
	result := syslogblocks.SyslogConfiguration{}
	result.ADDRTCP = "0.0.0.0:5141"
	result.SEVERITYLEVEL = 7
	return result
}
