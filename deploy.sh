#! /bin/bash

# This script is used to deploy the xBot application.
go build -o xbot ./cmd/xbot
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi