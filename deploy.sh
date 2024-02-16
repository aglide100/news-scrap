#!/bin/bash

docker stack deploy -c <(docker-compose -f app.yml config) scrap
