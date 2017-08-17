GO = go

all: build

help:
	@echo $(TAG)
	@echo "build     - go build"
	@echo "run       - go run"
	@echo "test      - curl predict api"
	@echo "loadtest  - load test predict api"

build: checkdeps
	$(GO) build -i $(GOFLAGS)
	$(GO) test -i ./...

test: checkdeps
	curl -s -H "Content-Type: application/json" -d @talk.json -X POST http://localhost:3002/predict?top=3 | jq

loadtest:
	$(GO) get github.com/rakyll/hey
	hey -n 5000 -c 10 -H "Content-Type: application/json" -D ./talk.json -m POST http://localhost:3002/predict

checkdeps:
	$(GO) get -d -v

run: build
	go run main.go

clean:
	$(GO) clean

.PHONY: build test loadtest run checkdeps clean
