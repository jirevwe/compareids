#!/bin/bash

# Small Scale Testing Script
echo "Starting small scale ID comparison tests..."

# Copy environment configuration
cp env.small .env

# Run the tests
docker compose up --build

echo "Small scale tests completed!"
echo "Results are available in the ./results directory"
