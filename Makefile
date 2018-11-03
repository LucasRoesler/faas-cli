.GIT_COMMIT=$(shell git rev-parse HEAD)
.GIT_VERSION=$(shell git describe --tags 2>/dev/null || echo "$(.GIT_COMMIT)")
.GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(.GIT_UNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif

GO_FILES?=$$(find . -name '*.go' |grep -v vendor)
TAG?=latest

.PHONY: build
build:
	./build.sh

.PHONY: build_redist
build_redist:
	./build_redist.sh

.PHONY: build_samples
build_samples:
	./build_samples.sh

.PHONY: local-fmt
local-fmt:
	gofmt -l -d $(GO_FILES)

.PHONY: local-goimports
local-goimports:
	goimports -w $(GO_FILES)

.PHONY: test-unit
test-unit:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /template/ | grep -v build) -cover

ci-armhf-push:
	(docker push openfaas/faas-cli:$(TAG)-armhf)
ci-armhf-build:
	(./build.sh $(TAG)-armhf)

.PHONY: test-templating
PORT?=38080
FUNCTION?=templating-test-func
FUNCTION_UP_TIMEOUT?=30
.EXPORT_ALL_VARIABLES:
test-templating:
	./build_integration_test.sh


install: | $(GOPATH)/bin/faas-cli

$(GOPATH)/bin/faas-cli: $(shell find . -type f -name '*.go') ## install the cli to your local system
	@echo "Installing faas-cli"
	go install --ldflags "-s -w \
        -X github.com/openfaas/faas-cli/version.GitCommit=$(.GIT_COMMIT) \
        -X github.com/openfaas/faas-cli/version.Version=$(.GIT_VERSION)" \
        -a -installsuffix cgo

	@echo "Installed:"
	@faas-cli version

