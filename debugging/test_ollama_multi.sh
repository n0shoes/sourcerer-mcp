#!/bin/bash
for i in {1..10}; do
  echo "Request $i:"
  response=$(curl -s -X POST http://localhost:11434/api/embed -d "{\"model\": \"nomic-embed-text\", \"input\": \"test embedding number $i with some text\"}")
  if echo "$response" | grep -q '"error"'; then
    echo "$response" | jq -r '.error'
  elif echo "$response" | grep -q '"embeddings"'; then
    echo "✓ Success"
  else
    echo "⚠️  Unexpected response: $response"
  fi
  sleep 0.2
done
