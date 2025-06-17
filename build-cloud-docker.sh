#!/bin/bash

if [ "$#" -eq 0 ]; then
    echo "No arguments provided. Available parameters:"
    echo "1. PREFIX (default: plan-pocker)"
    echo "2. TAG (default: latest, can be set to a Git tag auto if TAG=git)"
    echo "3. SUFFIX (default: '')"
    read -p "Are you sure you want to use the default values? (y/n): " answer
    if [ "$answer" != "y" ]; then
        echo "Script terminated."
        exit 1
    fi
fi

PREFIX=${1:-plan-pocker}
TAG=${2:-latest}
SUFFIX=${3:-""}

if [ -n "$SUFFIX" ]; then
    CONTAINER_NAME_SUFFIX="-$SUFFIX"
else
    CONTAINER_NAME_SUFFIX=""
fi

if [ "$TAG" == "git" ]; then
    if command -v git >/dev/null 2>&1; then
        COMMIT_HASH=$(git rev-parse --short HEAD)
        TAG=$COMMIT_HASH
        echo "Used git as tag: ${TAG}"
    else
        echo "Error: git is not available. Please ensure git is installed and accessible in PATH."
        exit 1
    fi
fi


docker build -t "${PREFIX}-websocket-server${POSTFIX}:${TAG}" -f ./cloud/websocket-server/Dockerfile .
docker build -t "${PREFIX}-frontend${POSTFIX}:${TAG}"         -f ./cloud/frontend/Dockerfile         .
docker build -t "${PREFIX}-game${POSTFIX}:${TAG}"             -f ./cloud/game/Dockerfile             .
docker build -t "${PREFIX}-api-server${POSTFIX}:${TAG}"       -f ./cloud/api-server/Dockerfile       .
docker build -t "${PREFIX}-nginx${POSTFIX}:${TAG}"            -f ./cloud/nginx/Dockerfile            ./cloud/nginx/


if [ "$TAG" != "latest" ]; then
    for SERVICE in websocket-server frontend game api-server nginx; do
        docker tag "${PREFIX}-${SERVICE}${CONTAINER_NAME_SUFFIX}:${TAG}" "${PREFIX}-${SERVICE}${CONTAINER_NAME_SUFFIX}:latest"
    done;
fi