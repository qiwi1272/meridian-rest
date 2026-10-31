.PHONY: rest examples test test-integration deps version check-version tag release all fmt tidy vet build

# Default CI / release checks: vet, all packages tests, build.
# Integration (live Testnet): MERIDIAN_INTEGRATION=1 — see `make test-integration` (runs go test ./...).

# Releases are manual: `make release NEW_VERSION=1.1.0` (or the Release workflow).
# The module path has no /vN suffix, so only v0/v1 versions are valid.
LATEST_TAG := $(shell git describe --tags --abbrev=0 --match 'v[01].*' 2>/dev/null || echo none)
RELEASE := $(patsubst v%,%,$(NEW_VERSION))

fmt:
	go fmt . ./...

tidy:
	go mod tidy

vet:
	go vet . ./...

build:
	go build ./...

rest: vet test
	go build ./...

# Optional: hits api.meridiantest.net (requires network; may need MERIDIAN_PK). Runs the full tree.
test:
	MERIDIAN_INTEGRATION=1 go test -v ./...

examples:
	mkdir -p bin
	go build -o bin/example_account_balance examples/account_balance/main.go
	go build -o bin/example_limit_multiple examples/limit_multiple/main.go
	go build -o bin/example_limit_single examples/limit_single/main.go
	go build -o bin/example_positions examples/positions/main.go
	go build -o bin/example_cancel_replace examples/cancel_replace/main.go
	go build -o bin/example_twap examples/twap/main.go
	go build -o bin/example_chase examples/chase/main.go

deps:
	go get -u ./...

version:
	@echo "Latest v0/v1 tag: $(LATEST_TAG)"

check-version:
	@echo "$(RELEASE)" | grep -Eq '^[01]\.[0-9]+\.[0-9]+$$' || \
		{ echo "NEW_VERSION must be 0.x.y or 1.x.y, e.g. make release NEW_VERSION=1.1.0"; exit 1; }
	@! git rev-parse -q --verify "refs/tags/v$(RELEASE)" >/dev/null || \
		{ echo "tag v$(RELEASE) already exists"; exit 1; }

tag:
	git tag -a v$(RELEASE) -m "Release v$(RELEASE)"
	git push origin v$(RELEASE)

release: check-version rest tag
	gh release create v$(RELEASE) \
		--latest \
		--title "v$(RELEASE)" \
		--notes "Release v$(RELEASE)"

all: fmt tidy rest
