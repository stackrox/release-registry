#!/bin/bash

set -uo pipefail

golangci-lint run ./...
