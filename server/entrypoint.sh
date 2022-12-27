#!/bin/sh

if [ -z "$1" ]; then
    /app/main
else
    $*
fi

