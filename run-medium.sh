#!/bin/bash

# Medium Scale Testing Script
echo "Starting medium scale ID comparison tests..."

# Copy environment configuration
cp env.medium .env

# Run the tests
docker compose up --build

echo "Medium scale tests completed!"
echo "Results are available in the ./results directory"
