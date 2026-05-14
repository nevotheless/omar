VERSION ?= 2026.5.0-dev
VARIANT ?= omarchy

.PHONY: build test test-race test-cover test-short test-ci lint lint-all clean \
        image image-shell image-vm coverage

build:
	go build -ldflags="-X 'github.com/nevotheless/omar/internal/version.Version=$(VERSION)'" -o bin/omar ./cmd/omar

test:
	go test ./... -v

test-race:
	go test ./... -race -v

test-short:
	go test ./... -short -count=1

test-cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic -v
	go tool cover -func=coverage.out | tail -1

coverage: test-cover
	go tool cover -html=coverage.out -o coverage.html
	@echo "→ coverage.html generated — open in browser"

test-ci: test-race test-cover

lint:
	go vet ./...

lint-all:
	go vet ./...
	test -z "$$(gofmt -l .)" || (echo "Formatting issues:"; gofmt -l .; exit 1)

clean:
	rm -rf bin/ output/ output-*/ coverage.out coverage.html

image: images/mkosi.conf.d/packages.conf
	sudo mkosi -d arch --format=oci -C images --image-tag=v$(VERSION) build

image-shell:
	sudo mkosi -d arch -C images shell

image-vm:
	sudo mkosi -d arch -C images vm

images/mkosi.conf.d/packages.conf: images/generate-package-configs.sh images/packages/*.txt
	VARIANT=$(VARIANT) ./images/generate-package-configs.sh
