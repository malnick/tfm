#!/bin/bash
SOURCE_DIR=$(git rev-parse --show-toplevel)

function cleanup() {
	docker rmi -f deployd
}
trap cleanup INT HUP TERM KILL 

echo "building container for testing deployd"
echo "docker build -t deployd $SOURCE_DIR"
docker build -t deployd $SOURCE_DIR
docker run -i --name deployd --rm deployd make test
CODE=$?
cleanup
exit $CODE
