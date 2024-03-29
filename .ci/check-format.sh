#!/bin/sh
echo "Checking files for correct go formatting..."

WRONG_FORMAT=$(gofmt -l .)

# if gofmt found no files with wrong formatting, exit ok
[ -z "$WRONG_FORMAT" ] && exit 0

# else print and fail
echo >&2 "Unformatted files found!\n\n$WRONG_FORMAT\n\nPlease format them using 'go fmt ./...'!"
exit 1