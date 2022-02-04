#!/bin/bash

if [ ! -f ./config.yaml ]; then
    echo "Config file not founded!"
    exit 1
fi


go build -o ruller && ./ruller ./config.yaml $@