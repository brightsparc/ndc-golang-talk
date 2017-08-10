# Verify endpoint is succesful
curl -H "Content-Type: application/json" -d @talk.json -X POST http://localhost:3002/predict

# Load test endpoint
hey -n 10000 -c 100 -H "Content-Type: application/json" -D ./talk.json -m POST http://localhost:3002/predict
