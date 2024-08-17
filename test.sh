#!/bin/sh

# Run tests, check the formatting and check for suspicious code constructs using golang tools

set -o errexit
set -o nounset

echo "Running tests:"
go test ./...
echo

echo "Checking gofmt: "
ERRS=$(find . -type f -name "*.go" | xargs gofmt -l 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL - the following files need to be gofmt'ed:"
    for e in ${ERRS}; do
        echo "    $e"
    done
    echo
    exit 1
fi
echo "PASS"
echo

echo "Checking go vet: "
ERRS=$(go vet ./... || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"
echo
