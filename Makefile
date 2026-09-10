APP := kbot
REGISTRY := ghcr.io
OWNER := ekucher

VERSION ?= $(shell git describe --tags --abbrev=0)
COMMIT ?= $(shell git rev-parse --short=7 HEAD)

OS ?= linux
ARCH ?= amd64

TAG := $(VERSION)-$(COMMIT)
IMAGE := $(REGISTRY)/$(OWNER)/$(APP):$(TAG)-$(OS)-$(ARCH)

.PHONY: test build image push image-name image-tag clean

test:
	go test ./...

build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -trimpath -ldflags="-s -w" -o bin/$(APP) .

image:
	docker build --platform $(OS)/$(ARCH) --build-arg TARGETOS=$(OS) --build-arg TARGETARCH=$(ARCH) -t $(IMAGE) .

push:
	docker push $(IMAGE)

image-name:
	@echo $(IMAGE)

image-tag:
	@echo $(TAG)

clean:
	rm -rf bin
