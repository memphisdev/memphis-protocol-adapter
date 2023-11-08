package adapter

import (
	lgr "github.com/memphisdev/memphis-rest-gateway/logger"
)

type LoggerFactory func() *lgr.Logger
