#!/bin/sh

if [ -z "$1" ]; then
    echo "Starting server container ..."
    /app/main
else
    $*
fi
