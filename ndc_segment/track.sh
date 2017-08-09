# Make sure we have created the delivery stream
aws firehose --endpoint http://localhost:4573 \
  create-delivery-stream --delivery-stream-name=test-stream

# Send track event to API
curl -u 8GT90rhOugcEGAqdcBWP8muss7iSwwHy: \
  -H "Content-Type: application/json" -d @track.json \
  -X POST ${1-https://api.segment.io/v1/track}
