GO = GOGC=off go
GOFLAGS = -ldflags "-X main.version=$(VERSION)"
TAG = brightsparc/ndc_segment
VERSION = v1.0.1

all: build test

help:
	@echo $(TAG)
	@echo "build     - go build"
	@echo "docker    - docker build"
	@echo "run       - docker-compose up"
	@echo "publish   - docker push"
	@echo "test      - curl tack api"
	@echo "loadtest  - load test track api"
	@echo "clean     - remove temp files"

build: checkdeps
	$(GO) build -i $(GOFLAGS)
	$(GO) test -i ./...

test: checkdeps
	curl -H "Content-Type: application/json" -d @track.json -X POST http://localhost:3001/v1/track

loadtest: test
	$(GO) get -u github.com/rakyll/hey
	hey -n 100 -c 10 -H "Content-Type: application/json" -D ./track.json -m POST http://localhost:3001/v1/track

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
