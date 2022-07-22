#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo "Missing required arguments"
    echo "run as: ./release.sh APPLICATION_ID VERSION OUTPUT_FILE"
    exit 1
fi

CGO_ENABLED=0 go build -ldflags=\
"-X 'github.com/UnseenWizzard/pocket-cli/pkg/util.PocketAppId=$1'\
 -X 'github.com/UnseenWizzard/pocket-cli/pkg/util.Version=$2'\
 -w -s -extldflags '-static'"\
 -o $3