#!/bin/bash
logger  --rfc5424 --server 127.0.0.1 --port 5141 --priority user.alert  $(date +%d-%m-%Y_%H-%M-%S)