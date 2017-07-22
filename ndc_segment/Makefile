# do not specify a full path for go since travis will fail
VERSION = v1.0.1
TAG = brightsparc/ndc_segment
DOCKERHOST = dockerhost

GO = GOGC=off go
GOFLAGS = -ldflags "-X main.version=$(VERSION)"
GOVENDOR = $(shell which govendor)

all: build test

help:
	@echo $(TAG)
	@echo "build     - go build"
	@echo "install   - go install"
	@echo "test      - go test"
	@echo "gofmt     - go fmt"
	@echo "linux     - go build linux/amd64"
	@echo "clean     - remove temp files"

build: checkdeps
	$(GO) build -i $(GOFLAGS)
	$(GO) test -i ./...

test: checkdeps
	$(GO) test -v -test.timeout 15s `go list ./... | grep -v '/vendor/'`

checkdeps:
	[ -x "$(GOVENDOR)" ] || $(GO) get -u github.com/kardianos/govendor
	govendor list +e | grep '^ e ' && { echo "Found missing packages. Please run 'govendor add +e'"; exit 1; } || : echo

gofmt:
	gofmt -w `find . -type f -name '*.go' | grep -v vendor`

linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS)

docker: linux
	docker build -q -t $(TAG) .
	docker tag $(TAG) $(TAG):$(VERSION)

install:
	$(GO) install $(GOFLAGS)

firehose:
	# Run the localstack, create and list streams
	docker run -d -e SERVICES:firehose -p 4573:4573 -p 8080:8080 --name firehose atlassianlabs/localstack
	aws firehose --endpoint http://$(DOCKERHOST):4573/ create-delivery-stream --delivery-stream-name="test-stream"

clean:
	$(GO) clean

.PHONY: build linux gofmt install docker firehose test clean
