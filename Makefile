INTEGRATION  := $(shell basename $(shell pwd))
BINARY_NAME   = nr-$(INTEGRATION)
GO_PKGS      := $(shell go list ./... | grep -v "/vendor/")
GO_FILES     := $(shell find src -type f -name "*.go")
VALIDATE_DEPS = github.com/golang/lint/golint
DEPS          = github.com/kardianos/govendor
TEST_DEPS     = github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml

all: build

build: clean validate compile test

clean:        
	@echo "=== $(INTEGRATION) === [ clean ]: removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

validate-deps:
	@echo "=== $(INTEGRATION) === [ validate-deps ]: installing validation dependencies..."
	@go get -v $(VALIDATE_DEPS)

validate-only:
	@printf "=== $(INTEGRATION) === [ validate ]: running gofmt... "
# `gofmt` expects files instead of packages. `go fmt` works with
# packages, but forces -l -w flags.
	@OUTPUT="$(shell gofmt -l $(GO_FILES))" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Incorrect syntax in the following files:" ;\
		echo "$$OUTPUT" ;\
		exit 1 ;\
	fi
	@printf "=== $(INTEGRATION) === [ validate ]: running golint... "
	@OUTPUT="$(shell golint $(GO_PKGS))" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Issues found:" ;\
		echo "$$OUTPUT" ;\
		exit 1 ;\
	fi
	@printf "=== $(INTEGRATION) === [ validate ]: running go vet... "
	@OUTPUT="$(shell go vet $(GO_PKGS))" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Issues found:" ;\
		echo "$$OUTPUT" ;\
		exit 1;\
	fi

validate: validate-deps validate-only

compile-deps:
	@echo "=== $(INTEGRATION) === [ compile-deps ]: installing build dependencies..."
	@go get $(DEPS)
	@govendor sync

compile-only:
	@echo "=== $(INTEGRATION) === [ compile ]: building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

compile: compile-deps compile-only

test-deps: compile-deps
	@echo "=== $(INTEGRATION) === [ test-deps ]: installing testing dependencies..."
	@go get -v $(TEST_DEPS)

test-only:
	@echo "=== $(INTEGRATION) === [ test ]: running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml

test: test-deps test-only

.PHONY: all build clean validate-deps validate-only validate compile-deps compile-only compile test-deps test-only test
