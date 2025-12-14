.PHONY: default
default: build lint test

.PHONY: build
build:
	go build -o /dev/null .

.PHONY: test
test:
	go test -v ./...

GOLANGCI_LINT_VERSION=v2.7.2
GOLANGCI_LINT = $(shell go env GOBIN)/golangci-lint
$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

VERSION_FILE=version.go
VERSION = $(shell cat $(VERSION_FILE) | grep 'const Version' | cut -d '"' -f 2)

.PHONY: version
version:
	@echo $(VERSION)

VERSION_TAG_NAME = v$(VERSION)
.PHONY: tag_push
tag_push:
	git tag $(VERSION_TAG_NAME)
	git push origin $(VERSION_TAG_NAME)

TEXT_TEMPLATE_CLI=$(shell go env GOPATH)/bin/text-template-cli
$(TEXT_TEMPLATE_CLI):
	go install github.com/akm/text-template-cli@latest

README.md: README-gen

.PHONY: README-gen
README-gen:
	$(TEXT_TEMPLATE_CLI) README.md.tmpl > README.md
