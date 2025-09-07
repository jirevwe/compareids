#!/bin/bash

# Billion Scale Testing Script
echo "Starting billion scale ID comparison tests..."
echo "WARNING: This will take 24-72 hours and requires significant resources!"

# Confirm before proceeding
read -p "Are you sure you want to proceed with billion-scale testing? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Billion scale testing cancelled."
    exit 1
fi

# Copy environment configuration
cp env.billion .env

# Run the tests
docker compose up --build

echo "Billion scale tests completed!"
echo "Results are available in the ./results directory"
