#!/bin/bash
seq -w 1 1 1000000 > /tmp/syslog.1 &
seq -w 1000001 1 2000000 > /tmp/syslog.2 &
seq -w 2000001 1 3000000 > /tmp/syslog.3 &
seq -w 3000001 1 4000000 > /tmp/syslog.4 &

sleep 5

logger  --rfc5424 --tcp --server 127.0.0.1 --port 5141 --priority kern.emerg  -f /tmp/syslog.1 &
logger  --rfc3164 --udp --server 127.0.0.1 --port 5141 --priority user.crit  -f /tmp/syslog.2 &

