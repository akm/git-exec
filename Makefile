.PHONY: default
default: build test

.PHONY: build
build:
	go build -o /dev/null .

.PHONY: test
test:
	go test -v .
