#!/bin/bash

# Run linters on the project
golangci-lint run -c ./_lint/.production.golangci.json . && \
cd validation && go build -o validation.bin . && \
./validation.bin validate-printf && \
validation && ./validation.bin validate-errorf && \
./validation.bin newbench && \
rm validation.bin && \
cd .. && \
go test -v -race -coverpkg=. -coverprofile=coverage.txt -covermode=atomic ./... && \
gosec -no-fail -fmt sarif -out gosec-results.sarif -exclude-dir=xprint_test -exclude-dir=benchmark ./... && \
govulncheck ./... && \
go vet ./... 


