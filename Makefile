.PHONY: default
default: build lint test

.PHONY: build
build:
	go build -o /dev/null .

.PHONY: test
test:
	go test -v .

GOLANGCI_LINT_VERSION=v1.61.0
GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint
$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...
