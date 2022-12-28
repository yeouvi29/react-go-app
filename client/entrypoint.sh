#!/bin/sh

if [ -z "$1" ]; then
    npm start
else
    $*
fi
