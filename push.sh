#!/bin/bash
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

docker build -t "$IMAGEURL" -f cmd/scrap/Dockerfile .

docker push "$IMAGEURL"
