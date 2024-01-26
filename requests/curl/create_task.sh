#!/bin/bash

curl -X POST --location "http://localhost:8080/quote/task" \
    -H "Content-Type: application/json" \
    -d '{
          "base": "'$1'",
          "counter": "'$2'"
        }'