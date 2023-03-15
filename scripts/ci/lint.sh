#!/bin/bash

set -uo pipefail

# Golang code
golangci-lint run ./...

# Proto
buf lint
