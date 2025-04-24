#!/bin/bash

UNFORMATTED=$(gofmt -l .)

if [ -n "$UNFORMATTED" ]; then
  echo -e "\033[31m[✗] Go files need formatting:\033[0m"
  echo "$UNFORMATTED"
  exit 1
else
  echo -e "\033[32m[✓] go fmt passed!\033[0m"
fi