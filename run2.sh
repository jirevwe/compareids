#!/bin/bash

# Exit on error
set -e

# Variables
DB_CONTAINER_NAME="postgres_test_container"
DB_PORT=5432
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="testdb"
DOCKER_IMAGE="postgres:15-alpine"
QUERIES_DIR="samples" # Folder containing query files
MOUNTED_DIR="/queries" # Mount point inside the container

# Cleanup function
cleanup() {
  echo "Stopping and removing container..."
  docker rm -f $DB_CONTAINER_NAME &>/dev/null || true
}

# Trap exit to cleanup on script termination
trap cleanup EXIT

# Start PostgreSQL in Docker
echo "Starting PostgreSQL container..."
docker run -d --name $DB_CONTAINER_NAME -e POSTGRES_PASSWORD=$DB_PASSWORD -e POSTGRES_DB=$DB_NAME -p $DB_PORT:5432 -v $(pwd)/$QUERIES_DIR:$MOUNTED_DIR $DOCKER_IMAGE

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
while ! docker exec -it $DB_CONTAINER_NAME pg_isready -U $DB_USER &>/dev/null; do
  sleep 1
done

echo "PostgreSQL is ready."

# Process each SQL file in the queries directory
for sql_file in $QUERIES_DIR/*.sql; do
  [ -e "$sql_file" ] || continue # Skip if no .sql files found

  base_name=$(basename "$sql_file" .sql)
  output_file="$QUERIES_DIR/$base_name.json"

  IN=$(tr -d '\n' < "$sql_file")

  echo "Processing $sql_file..."
  {
      echo "["
      IFS=';' read -ra lines <<< "$IN"
      for line in "${lines[@]}"; do
          start_time=$(docker exec -i $DB_CONTAINER_NAME echo $(($(date +%s%N)/1000000)))
          result=$(docker exec -i $DB_CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -t -q <<< "${line}")
          end_time=$(docker exec -i $DB_CONTAINER_NAME echo $(($(date +%s%N)/1000000)))

          duration=$((end_time - start_time))

          echo "  {\"query\": \"$(echo $line | sed 's/"/\\"/g')\", \"duration_ms\": $duration, \"result\": \"$result\" },"
      done
      echo "]"
  } > "$output_file"

  echo "Results written to $output_file"
done
