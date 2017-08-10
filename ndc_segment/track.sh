# Verify endpoint is succesful
curl -H "Content-Type: application/json" -d @track.json -X POST ${1-https://api.segment.io/v1/t}

# Load test endpoint
hey -n 1000 -c 10 -H "Content-Type: application/json" -D ./track.json -m POST ${1-https://api.segment.io/v1/t}
