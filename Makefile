PKG := $(shell go list .)

GOFILES := $(shell find . -type f -name '*.go' -print)

.PHONY: check
check: fmt-check lint test

.PHONY: deps
deps:
	go mod tidy

.PHONY: dev-deps
dev-deps:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: lint
lint: dev-deps
	staticcheck ./...

.PHONY: test
test:
	go test ./...

.PHONY: fmt-check
fmt-check: dev-deps
	goimports -l -local $(PKG) $(GOFILES) | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

.PHONY: fmt
fmt: dev-deps
	goimports -w -local $(PKG) $(GOFILES)
