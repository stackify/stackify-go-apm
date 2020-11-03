#!/bin/bash
set -e

echo "Running Tests and exclude example and instrumentation directories."
go test $(go list ./... | grep -v -e example -e instrumentation) -coverprofile coverage.html
echo "Done running tests."

echo "To view coverage file"
echo "Run: $ go tool cover -html=coverage.html"
