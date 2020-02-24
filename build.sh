#!/bin/bash
cd /opt/docker-compose/prestrafe-bot

docker-compose down
docker image rm prestrafe/prestrafe-bot:latest

git pull -p
docker build . --no-cache --tag prestrafe/prestrafe-bot:latest

docker-compose up -d
