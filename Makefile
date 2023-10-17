GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)
GOFILES := $(shell find . -name "*.go")
TESTFOLDER := $(shell $(GO) list ./... | grep -E 'utils$$')
TESTTAGS ?= "-v"
DOCKER ?= docker

.PHONY: test
test:
	$(GO) test $(TESTTAGS) -covermode=count -coverprofile=profile.out sbom-utilities/utils

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: lint
lint:
	$(shell golangci-lint run ./...)

.PHONY: clean
clean:
	$(shell rm profile.out)
	$(shell rm -rf bin)
	$(shell rm tmp.out)
	$(shell rm -rf utils/*temp*)
	$(shell rm -rf output)
	$(shell rm coverage.out)
	$(shell rm sbom-utilities)

.PHONY: build
build:
	$(GO) build -o bin/sbom-utils

.PHONY: docker
docker:
	$(DOCKER) build --build-arg ARCH=arm64 --tag sbom-utilities-pipe:dev .

.PHONY: docker-amd64
docker-amd64:
	$(DOCKER) buildx build --platform linux/amd64 --tag sbom-utilities-pipe:dev .

.PHONY: docker-run
docker-run:
	$(DOCKER) run --rm -it sbom-utilities-pipe:dev

.PHONY: docker-lint
docker-lint:
	$(DOCKER) run --rm -it \
		-v "$(shell pwd)":/build \
		--workdir /build \
		hadolint/hadolint:v2.12.0-alpine hadolint Dockerfile*

.PHONY: markdown-lint
markdown-lint:
	$(DOCKER) run --rm -it \
		-v "$(shell pwd)":/build \
		--workdir /build \
		markdownlint/markdownlint:0.13.0 *.md
