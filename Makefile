VERSION ?= 2026.5.0-dev
VARIANT ?= omarchy

.PHONY: build test clean lint image image-shell image-vm

build:
	go build -ldflags="-X 'github.com/nevotheless/omar/internal/version.Version=$(VERSION)'" -o bin/omar ./cmd/omar

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf bin/ output/ output-*/

image: images/mkosi.conf.d/packages.conf
	sudo mkosi -d arch --format=oci -C images --image-tag=v$(VERSION) build

image-shell:
	sudo mkosi -d arch -C images shell

image-vm:
	sudo mkosi -d arch -C images vm

images/mkosi.conf.d/packages.conf: images/generate-package-configs.sh images/packages/*.txt
	VARIANT=$(VARIANT) ./images/generate-package-configs.sh
