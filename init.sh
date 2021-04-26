#!/bin/bash
set -eu

cd $(dirname $0)

go build

bin_name="cwlogs"

if [ ! -e "$(pwd)/${bin_name}" ]; then
    echo "Error: Go Binary File not found"
    exit 1
fi

ln -s $(pwd)/${bin_name} /usr/local/bin/cwlogs
chmod 755 /usr/local/bin/cwlogs
