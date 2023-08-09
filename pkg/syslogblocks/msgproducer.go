package syslogblocks

import "github.com/g41797/sputnik"

type MsgProducer interface {
	Connect(cf sputnik.ConfFactory) error
	Produce(msg sputnik.Msg) error
	Disconnect()
}
