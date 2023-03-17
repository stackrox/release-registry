#!/bin/bash

GOLANGCI_LINT_VERSION="1.51.2"
INSTALLED_GOLANGCI_LINT_VERSION=$(golangci-lint version --format short || echo "not installed")

if [ "$GOLANGCI_LINT_VERSION" != "${INSTALLED_GOLANGCI_LINT_VERSION}" ]; then
    curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/v${GOLANGCI_LINT_VERSION}/install.sh" \
        | sh -s -- -b "$(go env GOPATH)/bin" "v${GOLANGCI_LINT_VERSION}"
else
    echo "golangci-lint already in version ${GOLANGCI_LINT_VERSION}, skipping..."
fi
