module github.com/memphisdev/memphis-protocol-adapter

go 1.19

require (
	github.com/g41797/kissngoqueue v0.1.5
	github.com/g41797/sputnik v0.0.6
	github.com/tkanos/gonfig v1.0.3
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/tkanos/gonfig => /home/g41797/go/pkg/mod/github.com/tkanos/gonfig/

replace gopkg.in/mcuadros/go-syslog.v2 => /home/g41797/go/pkg/mod/github.com/g41797/go-syslog
