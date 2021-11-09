.PHONY: all
all: build test

.PHONY: setup
setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest

.PHONY: build
build:
	go build -o esctl cmd/esctl/main.go

.PHONY: test
test:
	go test -count=1 -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html coverage.out -o coverage.html
	@coverage="$$(go tool cover -func coverage.out | grep 'total:' | awk '{print int($$3)}')"; \
	echo "The overall coverage is $$coverage%. Look at coverage.html for details.";

.PHONY: fix
fix:
	$(GOPATH)/bin/golangci-lint run --fix

.PHONY: lint
lint:
	$(GOPATH)/bin/golangci-lint run
