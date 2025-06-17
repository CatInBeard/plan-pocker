#!/bin/bash

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <command> [envfile=./cloud/test.env] [stackname=plan-pocker]"
    echo ""
    echo "Commands:"
    echo "  up    [service_name]            - Start the services. Can be used without parameters or with [service_name]."
    echo "  down  [service_name]            - Stop the services. Can be used without parameters or with [service_name]."
    echo "  scale <service_name> <replicas> - Scale a service. Requires [service_name] and [replicas]."
    echo ""
    echo "Notes:"
    echo "  - If [envfile] is not specified, environment variables will not be exported."
    echo "  - [stackname] is optional and defaults to 'plan-pocker'."
    exit 1
fi

COMMAND=$1
SERVICE_NAME=$2
REPLICAS=$3
STACK_NAME="plan-pocker"

mkdir -p cloud/data/logs
mkdir -p cloud/data/mysql
mkdir -p cloud/data/nginx
mkdir -p cloud/data/nginx-frontend
mkdir -p cloud/data/redisinsight
mkdir -p cloud/data/valkey

for arg in "$@"; do
    if [[ $arg == envfile=* ]]; then
        ENV_FILE="${arg#envfile=}"
    elif [[ $arg == stackname=* ]]; then
        STACK_NAME="${arg#stackname=}"
    fi
done

if [ -n "$ENV_FILE" ]; then
    if [ ! -f "$ENV_FILE" ]; then
        echo "Environment file $ENV_FILE not found!"
        exit 1
    fi

    export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

if [ "$TAG" == "git" ]; then
    if command -v git >/dev/null 2>&1; then
        COMMIT_HASH=$(git rev-parse --short HEAD)
        export TAG=$COMMIT_HASH
        echo "Used git as tag: ${TAG}"
    else
        echo "Error: git is not available. Please ensure git is installed and accessible in PATH."
        exit 1
    fi
fi

case $COMMAND in
    up)
        if [ -z "$SERVICE_NAME" ]; then
            echo "Building and deploying all services..."
            docker stack deploy -c swarm.yaml $STACK_NAME
        else
            echo "Building and deploying service $SERVICE_NAME..."
            docker stack deploy -c swarm.yaml $SERVICE_NAME
        fi
        ;;
    down)
        if [ -z "$SERVICE_NAME" ]; then
            echo "Stopping and removing all services..."
            docker stack rm $STACK_NAME
        else
            echo "Stopping and removing service $SERVICE_NAME..."
            docker stack rm $SERVICE_NAME
        fi
        ;;
    scale)
        if [ -z "$REPLICAS" ]; then
            echo "For the scale command, you must specify the number of replicas."
            exit 1
        fi
        echo "Scaling service $SERVICE_NAME to $REPLICAS replicas..."
        docker service scale ${SERVICE_NAME}=$REPLICAS
        ;;
    *)
        echo "Unknown command: $COMMAND"
        echo "Valid commands: up, down, scale"
        exit 1
        ;;
esac

if [ $? -eq 0 ]; then
    echo "Operation completed."
else
    echo "An error occurred during the operation."
    exit 1
fi
