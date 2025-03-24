#!/bin/bash

# Usage:
#   ./script.sh [--full|--test-only|--validate-only|--lint-gosec|--test-file <path>|--test-func <function>]
#
# Modes:
#  --full: (or no argument) Run the full chain of commands.
#  --test-only: Run only the "go test" command.
#  --validate-only: Run only the validation block.
#  --lint-gosec: Run govulncheck and then execute "./lint.sh ."
#  --test-file <path>: Run only "go test" on the provided file or package path.
#  --test-func <functionname>: Run only "go test" on "./..." with the -run filter for the specified function.

if [ "$1" == "--test-only" ]; then
    echo "Running tests only..."
    go test -v -race -coverpkg=. -coverprofile=coverage.txt -covermode=atomic -vet=off ./...
    exit 0
elif [ "$1" == "--validate-only" ]; then
    echo "Running validation only..."
    cd validation && \
      go build -o validation.bin . && \
      ./validation.bin validate-printf && \
      ./validation.bin validate-errorf && \
      ./validation.bin newbench && \
      rm validation.bin && \
    cd ..
    exit 0
elif [ "$1" == "--lint-gosec" ]; then
    echo "Running govulncheck and lint.sh..."
    govulncheck ./... && \
    exec ./lint.sh .
    exit 0
elif [ "$1" == "--test-file" ]; then
    if [ -z "$2" ]; then
      echo "Error: --test-file requires a path argument."
      exit 1
    fi
    echo "Running tests on path: $2"
    go test -v -race -coverpkg=. -coverprofile=coverage.txt -covermode=atomic -vet=off "$2"
    exit 0
elif [ "$1" == "--test-func" ]; then
    if [ -z "$2" ]; then
      echo "Error: --test-func requires a function name argument."
      exit 1
    fi
    echo "Running tests for function: $2"
    go test -v -race -coverpkg=. -coverprofile=coverage.txt -covermode=atomic -vet=off ./... -run "$2"
    exit 0
elif [ -z "$1" ] || [ "$1" == "--full" ]; then
    echo "Running full checks..."
    golangci-lint run -c ./_lint/.production.golangci.json . && \
    cd validation && \
      go build -o validation.bin . && \
      ./validation.bin validate-printf && \
      ./validation.bin validate-errorf && \
      ./validation.bin newbench && \
      rm validation.bin && \
    cd .. && \
    go test -v -race -coverpkg=. -coverprofile=coverage.txt -covermode=atomic -vet=off ./... && \
    gosec -no-fail -fmt sarif -out gosec-results.sarif -exclude-dir=xprint_test -exclude-dir=benchmark ./... && \
    govulncheck ./... && \
    go vet ./...
    exit 0
else
    echo "Usage: $0 [--full|--test-only|--validate-only|--lint-gosec|--test-file <path>|--test-func <function>]"
    exit 1
fi
