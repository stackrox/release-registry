TAG=$(shell git describe --tags --abbrev=10 --long)
TAGGED=$(shell git tag --contains | head)
ifneq (,$(TAGGED))
	# We're tagged. Use the tag explicitly.
	VERSION := $(TAGGED)
else
	# We're on a dev/PR branch
	VERSION := $(TAG)
endif

.PHONY: tag
tag:
	@echo $(VERSION)

#########
# Proto #
#########
PROTOC_VERSION=21.12
PROTO_SRC_DIR="proto/api/v1"
PROTO_THIRD_PARTY_DIR="proto/third_party"
PROTO_DEST_DIR="generated"

# Support different OS to diff local and CI
ifeq ($(shell uname -s),Linux)
PROTOC_URL = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
endif
ifeq ($(shell uname -s),Darwin)
PROTOC_URL = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-osx-x86_64.zip
endif

.PHONY: proto-install
proto-install:
	@echo "Install protoc $(PROTOC_VERSION)"
	@mkdir -p $(GOPATH)/bin
	@curl $(PROTOC_URL) -sLo /tmp/protoc.zip
	@unzip -o -d /tmp /tmp/protoc.zip bin/protoc
	install /tmp/bin/protoc $(GOPATH)/bin/protoc
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: proto-generate
proto-generate:
	@mkdir -p ${PROTO_DEST_DIR}
	protoc \
		-I=${PROTO_SRC_DIR} \
		-I=${PROTO_THIRD_PARTY_DIR} \
		--go_out=${PROTO_DEST_DIR} \
		--go-grpc_out=${PROTO_DEST_DIR} \
		${PROTO_SRC_DIR}/*

#######
# CLI #
#######
.PHONY: cli-run-local
cli-run-local:
	@go run cmd/cli/main.go

.PHONY: build-cli
cli-binary:
	go build -o build/relreg-cli cmd/cli/main.go

##########
# Server #
##########
.PHONY: server-run-local
server-run-local:
	@go run cmd/server/main.go

.PHONY: server-binary
server-binary:
	go build -o build/relreg-server cmd/server/main.go

server-image:
	docker build . -t quay.io/rhacs-eng/release-artifacts:$(VERSION)
