.PHONY: all
all:
	@echo "ERROR: No target selected."
	@make help
	@exit 1

TAG:=$(shell git describe --tags)
LOCAL_VALUES_FILE=deploy/chart/release-registry/configuration/values-${ENVIRONMENT}.yaml

#################
# CI & Building #
#################
.PHONY: install-linters
install-linters: ## Install linters and setup environment
	@mkdir -p outputs
	@./scripts/ci/install-linters.sh
	@./scripts/install-buf.sh

.PHONY: format
format: ## Format code
	@./scripts/ci/go-format.sh

.PHONY: lint
lint: ## Lint code
	@./scripts/ci/lint.sh

.PHONY: server-binary
server-binary: ## Builds server binary
	@go build -o build/release-registry cmd/server/main.go

.PHONY: server-image
server-image: ## Builds server image
	@DOCKER_BUILDKIT=1 docker build . \
		-f image/Dockerfile \
		--secret id=npmrc,src=${HOME}/.npmrc \
		-t quay.io/rhacs-eng/release-registry:${TAG}

.PHONY: server-image-push
server-image-push: ## Pushes server image to registry
	@docker push quay.io/rhacs-eng/release-registry:${TAG}

###################
# Helm Deployment #
###################
.PHONY: server-helm-template
server-helm-template: pre-check ## Renders the chart with Helm for debugging
	@gcloud secrets versions access latest \
		--secret "release-registry-${ENVIRONMENT}" \
		--project stackrox-infra \
	| \
	helm template \
		release-registry \
		deploy/chart/release-registry \
		--debug \
		--namespace release-registry \
		--set image.tag="${TAG}" \
		--values -

.PHONY: server-helm-deploy
server-helm-deploy: pre-check ## Deploys the server with Helm
	@gcloud secrets versions access latest \
		--secret "release-registry-${ENVIRONMENT}" \
		--project stackrox-infra \
	| \
	helm upgrade \
		release-registry \
		deploy/chart/release-registry \
		--install \
		--create-namespace \
		--namespace release-registry \
		--set image.tag="${TAG}" \
		--values -

#########################
# Local Helm Deployment #
#########################
.PHONY: server-helm-upload-secrets
server-helm-upload-local-values: pre-check ## Upload secrets from local configuration
	@gcloud secrets versions add "release-registry-${ENVIRONMENT}" \
		--data-file="${LOCAL_VALUES_FILE}" \
		--project stackrox-infra


.PHONY: server-helm-download-secrets
server-helm-download-local-values: pre-check ## Downloads secrets into local configuration
	@mkdir -p "$(dir ${LOCAL_VALUES_FILE})"
	@gcloud secrets versions access latest \
		--secret "release-registry-${ENVIRONMENT}" \
		--project stackrox-infra \
	> "${LOCAL_VALUES_FILE}"

.PHONY: server-helm-deploy-local-values
server-helm-deploy-local-values: pre-check ## Deploys the server with Helm and local configuration values
	@cat "${LOCAL_VALUES_FILE}" \
	| \
	helm upgrade \
		release-registry \
		deploy/chart/release-registry \
		--install \
		--create-namespace \
		--namespace release-registry \
		--set image.tag="${TAG}" \
		--values -

#########
# Tests #
#########
.PHONY: tests-unit
tests-unit: ## Runs all unit tests without cache
	@go test ./pkg/... -count=1

.PHONY: tests-integration
tests-integration: ## Runs all integration tests without cache
	@go test -v ./tests/integration/... -count=1

.PHONY: tests-e2e
tests-e2e: ## Runs all e2e tests without cache
	@go test -v ./tests/e2e/... -count=1

########
# Misc #
########
.PHONY: tag
tag: ## Describes current tag
	@echo ${TAG}

.PHONY: pre-check
pre-check:
ifndef ENVIRONMENT
	$(error ENVIRONMENT is not defined)
endif

.PHONY: init-dev-environment
init-dev-environment: ## Initializes local development environment after first clone
	@./scripts/install-buf.sh

.PHONY: server-renew-cert
server-renew-cert: ## Renews the gRPC gateway certificate
	@./scripts/cert/renew.sh

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
