#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo "Missing required arguments"
    echo "run as: ./release.sh APPLICATION_ID VERSION OUTPUT_FILE"
fi

CGO_ENABLED=0 go build -ldflags=\
"-X 'riedmann.dev/pocket-cli/pkg/util.PocketAppId=$1'\
 -X 'riedmann.dev/pocket-cli/pkg/util.Version=$2'\
 -w -s -extldflags '-static'"\
 -o $3