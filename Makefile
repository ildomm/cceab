# See https://golangci-lint.run/usage/install/
LINTER_VERSION = v1.59.1

# Variables needed when building binaries
VERSION := $(shell grep -oE -m 1 '([0-9]+)\.([0-9]+)\.([0-9]+)' CHANGELOG.md )
GIT_SHA := $(shell git rev-parse HEAD )

# To be used for dependencies not installed with gomod
LOCAL_DEPS_INSTALL_LOCATION = /usr/local/bin

.PHONY: clean
clean:
	rm -rf build
	mkdir -p build

.PHONY: deps
deps:
	go env -w "GOPRIVATE=github.com/ildomm/*"
	go mod download

.PHONY: build
build: deps build-api_handler build-validator

.PHONY: build-api_handler
build-api_handler: deps
	# Build the http server job binary
	cd cmd/api_handler && \
		go build -ldflags="-X main.semVer=${VERSION} -X main.gitSha=${GIT_SHA}" \
        -o ../../build/api_handler

.PHONY: build-validator
build-validator: deps
	# Build the background job binary
	cd cmd/validator && \
		go build -ldflags="-X main.semVer=${VERSION} -X main.gitSha=${GIT_SHA}" \
        -o ../../build/validator

.PHONY: unit-test
unit-test: deps
	go test -tags=testing -count=1 ./...

.PHONY: lint-install
lint-install:
	[ -e ${LOCAL_DEPS_INSTALL_LOCATION}/golangci-lint ] || \
	wget -O- -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b ${LOCAL_DEPS_INSTALL_LOCATION} ${LINTER_VERSION}

.PHONY: lint
lint: deps lint-install
	golangci-lint run

.PHONY: coverage-report
coverage-report: clean deps
	go test -tags=testing ./... \
		-coverprofile=build/cover.out github.com/ildomm/cceab/...
	grep -vE 'main\.go|test_helpers' build/cover.out > build/cover.temp && mv build/cover.temp build/cover.out
	go tool cover -html=build/cover.out -o build/coverage.html
	echo "** Coverage is available in build/coverage.html **"

.PHONY: coverage-total
coverage-total: clean deps
	go test -tags=testing ./... \
		-coverprofile=build/cover.out github.com/ildomm/cceab/...
	grep -vE 'main\.go|test_helpers' build/cover.out > build/cover.temp && mv build/cover.temp build/cover.out
	go tool cover -func=build/cover.out | grep total

.PHONY: security-setup
security-setup:
	go install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: security-check
security-check:
	govulncheck ./...