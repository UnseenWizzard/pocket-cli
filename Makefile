.PHONY: default build check fmt analysis test test-ci coverage-report coverage-check golangci-lint

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
	@.ci/coverage-report.sh 0

coverage-check:
	@.ci/coverage-report.sh 80

test-opts=
test:
	@go install gotest.tools/gotestsum@latest
	@gotestsum ${test-opts} -- -cover ./...

test-ci: test-opts=--junitfile unit-test.xml	
test-ci: test

golangci-lint:
	@golangci-lint run -E gosec -E bodyclose -E cyclop -E exhaustive -E godox