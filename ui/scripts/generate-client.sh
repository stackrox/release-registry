#!/usr/bin/env bash

# Generates TypeScript client based on the Swagger 2.0 definitions
# Should be invoked from release-registry/ui (grandparent) directory

OPENAPI_GENERATOR_CLI_IMAGE_TAG="v6.5.0"
OPENAPI_GENERATOR_CLI_IMAGE="openapitools/openapi-generator-cli:${OPENAPI_GENERATOR_CLI_IMAGE_TAG}"

CLIENT_DIR="src/client"

set -x

# paths below are relative to the git root "release-registry"
SWAGGER_FILE="gen/openapiv2/proto/api/v1/release.swagger.json"
GENERATOR_OUTPUT_DIR="ui/${CLIENT_DIR}/release"

docker run --privileged --rm -v "${PWD}/..:/local" "${OPENAPI_GENERATOR_CLI_IMAGE}" generate \
  -i "/local/${SWAGGER_FILE}" \
  -g typescript-axios \
  --skip-validate-spec \
  -o "/local/${GENERATOR_OUTPUT_DIR}"

# paths below are relative to the git root "release-registry"
SWAGGER_FILE="gen/openapiv2/proto/api/v1/quality_milestone.swagger.json"
GENERATOR_OUTPUT_DIR="ui/${CLIENT_DIR}/quality_milestone"

docker run --privileged --rm -v "${PWD}/..:/local" "${OPENAPI_GENERATOR_CLI_IMAGE}" generate \
  -i "/local/${SWAGGER_FILE}" \
  -g typescript-axios \
  --skip-validate-spec \
  -o "/local/${GENERATOR_OUTPUT_DIR}"

yarn prettier --ignore-path "!${CLIENT_DIR}" --write "${CLIENT_DIR}"
