#!/bin/bash

# Large Scale Testing Script
echo "Starting large scale ID comparison tests..."

# Copy environment configuration
cp env.large .env

# Run the tests
docker compose up --build

echo "Large scale tests completed!"
echo "Results are available in the ./results directory"
