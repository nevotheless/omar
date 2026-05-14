VERSION ?= 2026.5.0-dev

.PHONY: build test clean lint image

build:
	go build -ldflags="-X 'github.com/nevotheless/omar/internal/version.Version=$(VERSION)'" -o bin/omar ./cmd/omar

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf bin/ output/

image:
	sudo mkosi -d arch --format=oci -C images --image-tag=v$(VERSION) build

image-shell:
	sudo mkosi -d arch -C images shell

image-vm:
	sudo mkosi -d arch -C images vm
