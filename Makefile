build:
	cd examples/simple && go build
	cd examples/multisig && go build

install:
	go get github.com/stretchr/testify

test:
	go test -count=1 -covermode=atomic -coverprofile=coverage.txt ./...

.PHONY: build install test