#!/bin/bash
timestamp=$(date +%d-%m-%Y_%H-%M-%S)
logger  --rfc3164 --server 127.0.0.1 --port 5141 --priority user.alert  $timestamp
logger  --rfc5424 --server 127.0.0.1 --port 5141 --priority user.alert  $timestamp