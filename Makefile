# Local entry points for the conformance suite. CI runs `lint`, `test` and
# `acceptance`, one step per target.
#
# Nothing here publishes or writes platform state: the repository is the
# delivery.

SHELL := /bin/sh
GO ?= go

.PHONY: all build lint test acceptance package fmt-check vet clean

all: build lint test

# build produces the run entry operators execute.
build:
	$(GO) build -o bin/kos-conformance ./cmd/kos-conformance

fmt-check:
	@unformatted=$$(gofmt -l $$(git ls-files '*.go')); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt required for:"; echo "$$unformatted"; exit 1; \
	fi

vet:
	$(GO) vet ./...

lint: fmt-check vet

test:
	$(GO) test ./...

# acceptance runs the four design scenarios with verbose output; a failure here
# names the scenario that broke.
acceptance:
	$(GO) test ./tests/... -count=1 -v

# package writes this product's extension artifact through the SDK packing path;
# it prints the artifact path and the digest a component install names.
package:
	$(GO) run ./cmd/kos-conformance-package -out dist

clean:
	rm -rf bin dist
