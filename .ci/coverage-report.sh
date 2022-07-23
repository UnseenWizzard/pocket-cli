#!/bin/sh
if [ "$#" -ne 1 ]; then
  echo "Missing required argument TARGET COVERAGE PERCENT"
  exit 1
fi
TARGET=$1

go test -coverprofile coverage.out ./...
go tool cover -html coverage.out -o coverage.html
COVERAGE=$(go tool cover -func coverage.out | grep total | sed -n "s/total:\s*(statements)\s*\(.*\)%$/\1/p")

echo "Total test coverage: $COVERAGE%"

if [ $(echo "$COVERAGE < $TARGET" | bc -l) -gt 0 ]; then
  echo "::error:: Test coverage $COVERAGE% is below target of $TARGET%" && exit 1
fi