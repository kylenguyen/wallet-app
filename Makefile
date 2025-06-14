.PHONY: 

quality: lint local-test

test-setup:
	@if ! command -v gotestsum > /dev/null 2>&1; then \
		echo >&2 "gotestsum is not installed. Installing..."; \
		go install gotest.tools/gotestsum@latest; \
	fi

local-test: test-setup
	gotestsum -- -v -count=1 ./...

lint-setup:
	@if ! command -v golangci-lint > /dev/null 2>&1; then \
		echo >&2 "golangci-lint is not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.7; \
	fi

lint: lint-setup
	golangci-lint run ./...
