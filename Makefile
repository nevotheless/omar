.PHONY: build test clean lint image

BINARY=omar

build:
	go build -o bin/$(BINARY) ./cmd/omar

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf bin/ output/

image:
	sudo mkosi -d arch --format=oci -C images build

image-shell:
	sudo mkosi -d arch -C images shell

image-vm:
	sudo mkosi -d arch -C images vm
