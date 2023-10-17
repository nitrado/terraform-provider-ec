#Cleanup

fmt:
	@echo "==> Formatting source"
	@gofmt -s -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@echo "==> Done"
.PHONY: fmt

#Test

test:
	@go test -cover -race ./...
.PHONY: test

#Lint

lint:
	@golangci-lint run --config=.golangci.yml ./...
.PHONY: lint

#Build

build:
	@goreleaser release --clean --snapshot
.PHONY: build

# Schema Generation

schema-gen:
	@go run ./internal/cmd/schema-gen/
.PHONY: schema-gen
