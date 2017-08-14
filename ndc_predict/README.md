# NDC Predict

Go API for predicting tags against conference talk content.

## Prerequisites

Following is required to build and deploy.

* [Go 1.8](https://blog.golang.org/go1.8) or higher

## Getting started

Make sure your go environment is setup.

### Build

Compile go code:

```
$ go get && go build
```

### Run

Run the api locally.

```
$ go run main.go
2017/08/10 13:57:50 Loading model ../ndc_scraper/model.bin
2017/08/10 13:57:53 Model loaded in 2.928797763s
2017/08/10 13:57:53 NDC Predict API 1.0.0 listening on: :3002
```

### Testing

Test the API with curl.

```
$ curl -s -H "Content-Type: application/json" -d @talk.json -X POST http://localhost:3002/predict?top=3 | jq
[
  {
    "prob": 0.4511719,
    "label": "cloud"
  },
  {
    "prob": 0.17773439,
    "label": "microsoft"
  },
  {
    "prob": 0.08789065,
    "label": "web"
  }
]
```

Load test with [hey](https://github.com/rakyll/hey) for 100 concurrent users with ~10ms latency.

```
$ hey -n 10000 -c 10 -H "Content-Type: application/json" -D ./talk.json -m POST http://localhost:3002/predict

Summary:
  Total:	1.7166 secs
  Slowest:	0.1217 secs
  Fastest:	0.0002 secs
  Average:	0.0166 secs
  Requests/sec:	5825.4925
  Total data:	490000 bytes
  Size/request:	49 bytes

Status code distribution:
  [200]	10000 responses

Response time histogram:
  0.000 [1]	|
  0.001 [7801]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.003 [1705]	|∎∎∎∎∎∎∎∎∎
  0.004 [317]	|∎∎
  0.005 [84]	|
  0.006 [30]	|
  0.007 [21]	|
  0.008 [24]	|
  0.010 [6]	|
  0.011 [10]	|
  0.012 [1]	|

Latency distribution:
  10% in 0.0004 secs
  25% in 0.0006 secs
  50% in 0.0009 secs
  75% in 0.0013 secs
  90% in 0.0020 secs
  95% in 0.0026 secs
  99% in 0.0048 secs

```

## Authors

* Julian Bright - [brightsparc](https://github.com/brightsparc/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
