# NDC Segment

Go API for event tracking.

## Prerequisites

Following is required to build and deploy.

* [Go 1.8](https://blog.golang.org/go1.8) or higher
* [Docker 17.05](https://docs.docker.com/engine/installation) or higher.

## Getting started

Make sure your go environment is setup and docker engine is running.

### Build

Compile go code:

```
$ make
GOGC=off go get -d -v
GOGC=off go build -i -ldflags "-X main.version=v1.0.1"
GOGC=off go test -i ./...
GOGC=off go test -v -test.timeout 15s `go list ./... | grep -v '/vendor/'`
?   	github.com/brightsparc/ndc-golang-talk/ndc_segment	[no test files]
```

### Deploy

Make docker image:

```
$ make docker
docker build -t brightsparc/ndc_segment .
Sending build context to Docker daemon  9.619MB
Step 1/12 : FROM golang:1.8 as builder
 ---> 6ce094895555
Step 2/12 : WORKDIR /go/src/github.com/brightsparc/ndc-golang-talk/ndc_segment
 ---> Using cache
 ---> cc5c7940631a
Step 3/12 : COPY config.json .
 ---> Using cache
 ---> 0aebfb9fe67b
Step 4/12 : COPY main.go .
 ---> a255180c89bb
Removing intermediate container 98bb5eabdc8d
Step 5/12 : RUN go get -d -v
 ---> Running in 03cc9bb87a25
github.com/brightsparc/segment (download)
github.com/aws/aws-sdk-go (download)
github.com/gorilla/mux (download)
github.com/prometheus/client_golang (download)
github.com/beorn7/perks (download)
github.com/golang/protobuf (download)
github.com/prometheus/client_model (download)
github.com/prometheus/common (download)
github.com/matttproud/golang_protobuf_extensions (download)
github.com/prometheus/procfs (download)
github.com/segmentio/backo-go (download)
github.com/xtgo/uuid (download)
 ---> 944900c8f11b
Removing intermediate container 03cc9bb87a25
Step 6/12 : RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ndc_segment .
 ---> Running in 09376bcfb725
 ---> 532521ed1620
Removing intermediate container 09376bcfb725
Step 7/12 : FROM alpine:latest
 ---> 7328f6f8b418
Step 8/12 : RUN apk --no-cache add ca-certificates
 ---> Using cache
 ---> c05ee4a93054
Step 9/12 : WORKDIR /root/
 ---> Using cache
 ---> 8070532851ba
Step 10/12 : COPY --from=builder /go/src/github.com/brightsparc/ndc-golang-talk/ndc_segment .
 ---> f2f784591407
Removing intermediate container 81a922770f07
Step 11/12 : EXPOSE 3001
 ---> Running in 8ef7d5dc689b
 ---> a0c7493d5f6d
Removing intermediate container 8ef7d5dc689b
Step 12/12 : CMD ./ndc_segment
 ---> Running in 706dd38d5acc
 ---> 052a3a5bd978
Removing intermediate container 706dd38d5acc
Successfully built 052a3a5bd978
Successfully tagged brightsparc/ndc_segment:latest
docker tag brightsparc/ndc_segment brightsparc/ndc_segment:v1.0.1
```

### Run

Set `AWS` environment variables which will be passed through to docker container.

```
$ export AWS_REGION=us-west-2
$ export AWS_ACCESS_KEY_ID=xxx
$ export AWS_SECRET_ACCESS_KEY=yyy
```

Then run the the docker container with `docker-compose`:

```
$ make run
docker-compose up
ndcsegment_localstack_1 is up-to-date
Recreating ndcsegment_ndc_segment_1 ...
Recreating ndcsegment_ndc_segment_1 ... done
Recreating ndcsegment_prometheus_1 ...
Recreating ndcsegment_prometheus_1 ... done
Attaching to ndcsegment_localstack_1, ndcsegment_ndc_segment_1, ndcsegment_prometheus_1
ndc_segment_1  | 2017/08/10 03:49:15 Adding Segment handlers
ndc_segment_1  | 2017/08/10 03:49:15 NDC Segment API 1.0.0 listening on: :3001
ndc_segment_1  | 2017/08/10 03:49:15 Started forwarder processing
ndc_segment_1  | 2017/08/10 03:49:15 Delivery connecting to http://localstack:4573...
ndc_segment_1  | 2017/08/10 03:49:15 Created stream: arn:aws:firehose:us-east-1:000000000000:deliverystream/test-stream
ndc_segment_1  | 2017/08/10 03:49:15 Starting delivery processing
localstack_1   | /usr/lib/python2.7/site-packages/supervisor/options.py:296: UserWarning: Supervisord is running as root and it is searching for its configuration file in default locations (including its current working directory); you probably want to specify a "-c" argument specifying an absolute path to a configuration file for improved security.
localstack_1   |   'Supervisord is running as root and it is searching '
localstack_1   | 2017-08-10 02:07:49,781 CRIT Supervisor running as root (no user in config file)
localstack_1   | 2017-08-10 02:07:49,786 INFO supervisord started with pid 1
localstack_1   | 2017-08-10 02:07:50,791 INFO spawned: 'dashboard' with pid 8
localstack_1   | 2017-08-10 02:07:50,795 INFO spawned: 'infra' with pid 9
localstack_1   | (. .venv/bin/activate; bin/localstack web --port=8080)
localstack_1   | (. .venv/bin/activate; bin/localstack start)
localstack_1   | Starting local dev environment. CTRL-C to quit.
localstack_1   |  * Running on http://0.0.0.0:8080/ (Press CTRL+C to quit)
localstack_1   |  * Restarting with stat
localstack_1   |  * Debugger is active!
localstack_1   |  * Debugger PIN: 134-266-190
localstack_1   | 2017-08-10 02:07:52,689 INFO success: dashboard entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
localstack_1   | 2017-08-10 02:07:52,690 INFO success: infra entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
localstack_1   | Starting mock Firehose service (http port 4573)...
localstack_1   | Ready.
localstack_1   | /usr/lib/python2.7/site-packages/supervisor/options.py:296: UserWarning: Supervisord is running as root and it is searching for its configuration file in default locations (including its current working directory); you probably want to specify a "-c" argument specifying an absolute path to a configuration file for improved security.
localstack_1   |   'Supervisord is running as root and it is searching '
localstack_1   | 2017-08-10 03:46:38,740 CRIT Supervisor running as root (no user in config file)
localstack_1   | 2017-08-10 03:46:38,746 INFO supervisord started with pid 1
localstack_1   | 2017-08-10 03:46:39,750 INFO spawned: 'dashboard' with pid 9
localstack_1   | 2017-08-10 03:46:39,754 INFO spawned: 'infra' with pid 10
localstack_1   | (. .venv/bin/activate; bin/localstack start)
localstack_1   | (. .venv/bin/activate; bin/localstack web --port=8080)
localstack_1   | Starting local dev environment. CTRL-C to quit.
localstack_1   |  * Running on http://0.0.0.0:8080/ (Press CTRL+C to quit)
localstack_1   |  * Restarting with stat
localstack_1   |  * Debugger is active!
localstack_1   |  * Debugger PIN: 134-266-190
localstack_1   | 2017-08-10 03:46:41,684 INFO success: dashboard entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
localstack_1   | 2017-08-10 03:46:41,684 INFO success: infra entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
localstack_1   | Starting mock Firehose service (http port 4573)...
localstack_1   | Ready.
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Starting prometheus (version=1.7.1, branch=master, revision=3afb3fffa3a29c3de865e1172fb740442e9d0133)" source="main.go:88"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Build context (go=go1.8.3, user=root@0aa1b7fc430d, date=20170612-11:44:05)" source="main.go:89"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Host details (Linux 4.9.36-moby #1 SMP Wed Jul 12 15:29:07 UTC 2017 x86_64 4cede5896380 (none))" source="main.go:90"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Loading configuration file /etc/prometheus/prometheus.yml" source="main.go:252"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Loading series map and head chunks..." source="storage.go:428"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="89 series loaded." source="storage.go:439"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Starting target manager..." source="targetmanager.go:63"
prometheus_1   | time="2017-08-10T03:49:18Z" level=info msg="Listening on :9090" source="web.go:259"
```

### Testing

Test the API with curl and [hey](https://github.com/rakyll/hey) with 1000 requests.

```
$ make loadtest
GOGC=off go get -d -v
curl -H "Content-Type: application/json" -d @track.json -X POST http://localhost:3001/v1/track
{ "success": true }
GOGC=off go get -u github.com/rakyll/hey
hey -n 100 -c 10 -H "Content-Type: application/json" -D ./track.json -m POST http://localhost:3001/v1/track

Summary:
  Total:	0.1284 secs
  Slowest:	0.0392 secs
  Fastest:	0.0010 secs
  Average:	0.0117 secs
  Requests/sec:	778.7686
  Total data:	1900 bytes
  Size/request:	19 bytes

Detailed Report:

	DNS+dialup:
  		Average:	0.0006 secs
  		Fastest:	0.0000 secs
  		Slowest:	0.0079 secs

	DNS-lookup:
  		Average:	0.0005 secs
  		Fastest:	0.0000 secs
  		Slowest:	0.0060 secs

	Request Write:
  		Average:	0.0001 secs
  		Fastest:	0.0000 secs
  		Slowest:	0.0010 secs

	Response Wait:
  		Average:	0.0109 secs
  		Fastest:	0.0009 secs
  		Slowest:	0.0391 secs

	Response Read:
  		Average:	0.0000 secs
  		Fastest:	0.0000 secs
  		Slowest:	0.0002 secs

Status code distribution:
  [200]	100 responses

Response time histogram:
  0.001 [1]	|∎
  0.005 [28]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.009 [21]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.012 [11]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.016 [15]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.020 [5]	|∎∎∎∎∎∎∎
  0.024 [7]	|∎∎∎∎∎∎∎∎∎∎
  0.028 [4]	|∎∎∎∎∎∎
  0.032 [2]	|∎∎∎
  0.035 [2]	|∎∎∎
  0.039 [4]	|∎∎∎∎∎∎

Latency distribution:
  10% in 0.0023 secs
  25% in 0.0044 secs
  50% in 0.0088 secs
  75% in 0.0162 secs
  90% in 0.0269 secs
  95% in 0.0352 secs
  99% in 0.0392 secs
```

## Authors

* Julian Bright - [brightsparc](https://github.com/brightsparc/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
