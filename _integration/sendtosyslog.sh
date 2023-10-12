#!/bin/bash
seq -w 1 1 1000000 > /tmp/syslog.5424
seq -w 1000001 1 2000000 > /tmp/syslog.3164
logger  --rfc5424 --tcp --server 127.0.0.1 --port 5141 --priority user.alert  -f /tmp/syslog.5424 &
logger  --rfc3164 --udp --server 127.0.0.1 --port 5141 --priority user.alert  -f /tmp/syslog.3164 &

