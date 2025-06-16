#!/bin/sh
docker build -t plan-pocker-websocket-server -f ./cloud/websocket-server/Dockerfile .
docker build -t plan-pocker-frontend         -f ./cloud/frontend/Dockerfile         .
docker build -t plan-pocker-game             -f ./cloud/game/Dockerfile             .
docker build -t plan-pocker-api-server       -f ./cloud/api-server/Dockerfile       .