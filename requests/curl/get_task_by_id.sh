#!/bin/bash
curl -X GET --location "http://localhost:8080/quote/task/'$1'" \
    -H "Content-Type: application/json"