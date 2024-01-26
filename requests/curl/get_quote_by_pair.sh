#!/bin/bash
curl -X GET --location "http://localhost:8080/quote?baseCurrency='$1'&counterCurrency='$2'" \
    -H "Content-Type: application/json"