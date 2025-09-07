#!/bin/bash

# Run tests with monitoring enabled
echo "Starting ID comparison tests with monitoring..."

# Check if scenario is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <scenario> [env_file]"
    echo "Scenarios: small, medium, large, billion"
    echo "Example: $0 medium env.medium"
    exit 1
fi

SCENARIO=$1
ENV_FILE=${2:-"env.$SCENARIO"}

# Check if environment file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Environment file $ENV_FILE not found!"
    exit 1
fi

echo "Using scenario: $SCENARIO"
echo "Using environment file: $ENV_FILE"

# Copy environment configuration
cp "$ENV_FILE" .env

# Run with monitoring profile
docker compose --profile monitoring up --build

echo "Tests with monitoring completed!"
echo "Results are available in the ./results directory"
echo "Monitoring data is available at http://localhost:9090"
