#!/bin/bash
set -e

echo "Running end-to-end tests..."
make build
make migrate-up
./main & PID=$$!; \
  echo "Started main with PID $$PID"; \
  sleep 2; \
  go clean -testcache && go test -v ./tests/e2e/...; \
  echo "Killing PID $$PID"; \
  kill $$PID || echo "Could not kill process";
make clean
