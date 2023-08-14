module github.com/memphisdev/memphis-protocol-adapter

go 1.19

require (
	github.com/RackSec/srslog v0.0.0-20180709174129-a4725f04ec91
	github.com/g41797/kissngoqueue v0.1.5
	github.com/g41797/sputnik v0.0.7
	github.com/nats-io/nats.go v1.28.0
	github.com/tkanos/gonfig v1.0.3
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/nats-io/nats-server/v2 v2.9.21 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/tkanos/gonfig => /home/g41797/go/pkg/mod/github.com/tkanos/gonfig/

replace gopkg.in/mcuadros/go-syslog.v2 => /home/g41797/go/pkg/mod/github.com/g41797/go-syslog
