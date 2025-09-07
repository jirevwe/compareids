#!/bin/bash

set -e

# Configuration
TEST_SCENARIO=${TEST_SCENARIO:-medium}
MAX_ROWS=${MAX_ROWS:-1000000}
PARALLEL_TESTS=${PARALLEL_TESTS:-1}
RESULTS_DIR=${RESULTS_DIR:-/app/results}
CONFIG_DIR=${CONFIG_DIR:-/app/config}

echo "Starting ID comparison tests with scenario: $TEST_SCENARIO"
echo "Max rows: $MAX_ROWS"
echo "Results directory: $RESULTS_DIR"

# Function to wait for PostgreSQL to be ready
wait_for_postgres() {
    echo "Waiting for PostgreSQL to be ready..."
    until pg_isready -h postgres -p 5432 -U postgres; do
        echo "PostgreSQL is unavailable - sleeping"
        sleep 2
    done
    echo "PostgreSQL is ready!"
}

# Function to install pg-ulid extension
install_pg_ulid() {
    echo "Installing pg-ulid extension..."
    if [ ! -d "pg-ulid" ]; then
        git clone --depth 1 https://github.com/andrielfn/pg-ulid.git
    fi
    cd pg-ulid
    make install
    cd ..
    echo "pg-ulid extension installed successfully"
}

# Function to run tests based on configuration
run_tests() {
    local config_file="$CONFIG_DIR/${TEST_SCENARIO}.json"
    
    if [ ! -f "$config_file" ]; then
        echo "Configuration file $config_file not found. Using default settings."
        echo "Running all tests with default row counts..."
        compareids all --host postgres --port 5432 --user postgres --password postgres --dbname postgres
        return
    fi
    
    echo "Using configuration from: $config_file"
    
    # Extract row counts and ID types from config
    local row_counts=$(jq -r '.row_counts[]' "$config_file" | tr '\n' ' ')
    local id_types=$(jq -r '.id_types[]' "$config_file" | tr '\n' ' ')
    
    echo "Row counts: $row_counts"
    echo "ID types: $id_types"
    
    # Run tests for each ID type and row count combination
    for id_type in $id_types; do
        for count in $row_counts; do
            # Skip if count exceeds MAX_ROWS
            if [ "$count" -gt "$MAX_ROWS" ]; then
                echo "Skipping $id_type with $count rows (exceeds MAX_ROWS: $MAX_ROWS)"
                continue
            fi
            
            echo "Running test for $id_type with $count rows..."
            compareids id "$id_type" --count "$count" --host postgres --port 5432 --user postgres --password postgres --dbname postgres
            
            # Small delay between tests to allow system to stabilize
            sleep 2
        done
    done
    
    # Merge results
    echo "Merging test results..."
    compareids merge --host postgres --port 5432 --user postgres --password postgres --dbname postgres
}

# Function to display system information
display_system_info() {
    echo "=== System Information ==="
    echo "Hostname: $(hostname)"
    echo "CPU cores: $(nproc)"
    echo "Memory: $(free -h | grep '^Mem:' | awk '{print $2}')"
    echo "Disk space: $(df -h / | tail -1 | awk '{print $4}')"
    echo "=========================="
}

# Function to monitor system resources
monitor_resources() {
    echo "Starting resource monitoring..."
    (
        while true; do
            echo "$(date): CPU: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)%, Memory: $(free | grep Mem | awk '{printf "%.1f%%", $3/$2 * 100.0}')"
            sleep 30
        done
    ) &
    MONITOR_PID=$!
}

# Main execution
main() {
    display_system_info
    
    # Wait for PostgreSQL
    wait_for_postgres
    
    # Install pg-ulid extension
    install_pg_ulid
    
    # Start resource monitoring
    monitor_resources
    
    # Run tests
    run_tests
    
    # Stop monitoring
    if [ ! -z "$MONITOR_PID" ]; then
        kill $MONITOR_PID 2>/dev/null || true
    fi
    
    echo "All tests completed successfully!"
    echo "Results are available in: $RESULTS_DIR"
    
    # Display summary
    if [ -f "$RESULTS_DIR/data.json" ]; then
        echo "=== Test Summary ==="
        jq -r '.IDTypes[]' "$RESULTS_DIR/data.json" | wc -l | xargs echo "Total ID types tested:"
        jq -r '.RowCounts[]' "$RESULTS_DIR/data.json" | wc -l | xargs echo "Total row count variations:"
        echo "==================="
    fi
}

# Handle signals for graceful shutdown
trap 'echo "Received signal, shutting down..."; exit 0' SIGTERM SIGINT

# Run main function
main "$@"
