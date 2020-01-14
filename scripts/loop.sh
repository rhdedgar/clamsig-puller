#!/bin/bash

if [ $# -lt 2 ] ; then
    echo "Usage: $(basename $0) <sleep_time> <command> [arg] ..."
    exit 1
fi

SLEEP_TIME=$1
shift

# The purpose of this script is to run something in an infinite loop
while true ; do
    "$@"
    sleep $SLEEP_TIME
done
