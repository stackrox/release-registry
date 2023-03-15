#!/bin/bash

set -uo pipefail

go fmt ./...
buf format -w
