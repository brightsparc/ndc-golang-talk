# Send track event to API
curl -H "Content-Type: application/json" -d @track.json -X POST ${1-https://api.segment.io/v1/t}

# send multiple requests
hey -n 1000 -c 10 -H "Content-Type: application/json" -D ./track.json -m POST ${1-https://api.segment.io/v1/t}
