#!/bin/env bash

API_DOCS_DIR="./pkg/client"

generate_client_docker() {
  docker run --rm -v "${PWD}":/local -u "$(id -u):$(id -g)" openapitools/openapi-generator-cli generate \
    -i /local/docs/swagger.yaml \
    -o /local/$API_DOCS_DIR \
    -g go \
    --package-name client

  exit_on_error "Could not create the client"
}

exit_on_error() {
  if [ $? -ne 0 ]; then
    echo "❌ $1"
    exit 1
  fi
}

generate_client_docker
