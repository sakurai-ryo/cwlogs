#!/bin/bash
set -eu

cd $(dirname $0)

FILE="/usr/local/bin/cwlogs"

if [ -e ${FILE} ]; then
    rm ${FILE}
fi

go build

bin_name="cwlogs"

if [ ! -e "$(pwd)/${bin_name}" ]; then
    echo "Error: Go Binary File not found"
    exit 1
fi

ln -s $(pwd)/${bin_name} ${FILE}
chmod 755 ${FILE}
