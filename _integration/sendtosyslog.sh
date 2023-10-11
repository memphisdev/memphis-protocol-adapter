#!/bin/bash
echo -n "" > /tmp/syslog.msdgs
seq -w 1 1 1000000 > /tmp/syslog.msdgs
logger  --rfc5424 --server 127.0.0.1 --port 5141 --priority user.alert  -f /tmp/syslog.msdgs
