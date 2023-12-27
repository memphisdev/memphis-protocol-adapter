package adapter

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"

	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-rest-gateway/conf"
	lgr "github.com/memphisdev/memphis-rest-gateway/logger"
	mconnector "github.com/memphisdev/memphis-rest-gateway/memphisSingleton"
	"github.com/nats-io/nats.go"
)

const (
	labelLen = 3
)

const (
	sourceName         = "prtcl-adptr"
	syslogsStreamName  = "$memphis_syslogs"
	syslogsInfoSubject = "extern.info"
	syslogsWarnSubject = "extern.warn"
	syslogsErrSubject  = "extern.err"
)

var _ sputnik.ServerConnector = &BrokerConnector{}
var _ io.Writer = &BrokerConnector{}

const connectorConfName = "connector"

type BrokerConnector struct {
	io.Writer
	conf       conf.Configuration
	nc         *nats.Conn
	l          atomic.Pointer[lgr.Logger]
	flags      int
	pidPrefix  string
	labelStart int
	baseSubj   string
	lblToSubj  map[string]string
}

func (c *BrokerConnector) Connect(cf sputnik.ConfFactory) (conn sputnik.ServerConnection, err error) {
	if c.IsConnected() {
		return c.getInfo, nil
	}

	var conf conf.Configuration

	if err = cf(connectorConfName, &conf); err != nil {
		return nil, err
	}

	return c.ConnectWithConfig(conf)
}

func (c *BrokerConnector) ConnectWithConfig(conf conf.Configuration) (conn sputnik.ServerConnection, err error) {

	c.conf = conf

	nc, err := mconnector.Connect(conf)

	if err != nil {
		return nil, err
	}

	c.nc = nc

	c.createLogger()

	return c.getInfo, nil
}

func (c *BrokerConnector) IsConnected() bool {
	if c == nil {
		return false
	}

	if c.nc == nil {
		return false
	}

	if !c.nc.IsConnected() {
		if !c.nc.IsClosed() {
			c.nc.Close()
		}
		c.l.Store(nil)
		return false
	}

	return true
}

func (c *BrokerConnector) Disconnect() {
	if !c.IsConnected() {
		return
	}
	c.l.Store(nil)
	c.nc.Close()
	return
}

func (c *BrokerConnector) getInfo() (*lgr.Logger, *nats.Conn) {
	if c == nil {
		return nil, nil
	}
	return c.l.Load(), c.nc
}

func (c *BrokerConnector) createLogger() {
	// From rest-gateway implementation
	c.flags = log.LstdFlags | log.Lmicroseconds
	c.pidPrefix = fmt.Sprintf("[%d] ", os.Getpid())
	c.labelStart = len(c.pidPrefix) + 28 //???
	c.baseSubj = fmt.Sprintf("%s.%s.", syslogsStreamName, sourceName)
	c.lblToSubj = map[string]string{
		"INF": syslogsInfoSubject,
		"WRN": syslogsWarnSubject,
		"ERR": syslogsErrSubject,
	}

	c.l.Store(lgr.NewLogger(log.New(c, c.pidPrefix, c.flags)))

	return
}

func (c *BrokerConnector) Write(p []byte) (n int, err error) {
	if !c.IsConnected() {
		return 0, errors.New("not connected")
	}

	if c.conf.CLOUD_ENV {
		return len(p), nil
	}

	label := string(p[c.labelStart : c.labelStart+labelLen])
	subjectSuffix, ok := c.lblToSubj[label]
	if !ok { // skip other labels
		return 0, nil
	}

	subject := c.baseSubj + subjectSuffix

	if err := c.nc.Publish(subject, p); err != nil {
		return 0, err
	}

	return len(p), nil
}

func PrepareTLS(CLIENT_CERT_PATH, CLIENT_KEY_PATH, ROOT_CA_PATH string) (*tls.Config, error) {

	if CLIENT_CERT_PATH == "" || CLIENT_KEY_PATH != "" || ROOT_CA_PATH != "" {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(CLIENT_CERT_PATH, CLIENT_KEY_PATH)
	if err != nil {
		return nil, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	TLSConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	TLSConfig.Certificates = []tls.Certificate{cert}
	certs := x509.NewCertPool()

	pemData, err := os.ReadFile(ROOT_CA_PATH)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	TLSConfig.RootCAs = certs

	return TLSConfig, nil
}
