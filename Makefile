.PHONY: build check fmt analysis test test-ci coverage-report

default: check

build: 
	@go build -v ./...

check: build analysis test

fmt: 
	@go fmt ./...

analysis:
	@.ci/check-format.sh
	@go vet ./...

coverage-report:
	@go test -coverprofile coverage.out ./...
	@go tool cover -html coverage.out -o coverage.html

test-opts=
test:
	@go install gotest.tools/gotestsum@latest
	@gotestsum ${test-opts} -- -cover ./...

test-ci: test-opts=--junitfile unit-test.xml	
test-ci: test