.PHONY: lint vet staticcheck gosec

# Lint the Go code
lint:
	@echo "Running golint..."
	go run golang.org/x/lint/golint@latest ./...

# Run gosec to check for security vulnerabilities
gosec:
	@echo "Running gosec..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
