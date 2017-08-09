GO = GOGC=off go
GOFLAGS = -ldflags "-X main.version=$(VERSION)"
TAG = brightsparc/ndc_segment
VERSION = v1.0.1

all: build test

help:
	@echo $(TAG)
	@echo "build     - go build"
	@echo "test      - go test"
	@echo "docker    - docker build"
	@echo "run       - docker-compose up"
	@echo "publish   - docker push"
	@echo "clean     - remove temp files"

build: checkdeps
	$(GO) build -i $(GOFLAGS)
	$(GO) test -i ./...

test: checkdeps
	$(GO) test -v -test.timeout 15s `go list ./... | grep -v '/vendor/'`

checkdeps:
	$(GO) get -d -v

docker:
	docker build -t $(TAG) .
	docker tag $(TAG) $(TAG):$(VERSION)

run:
	docker-compose up

publish: docker
	docker push $(TAG):$(VERSION)

clean:
	$(GO) clean

.PHONY: build test docker run publish clean
