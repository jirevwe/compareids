package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config for database connection
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

// TestResult holds the result of a test
type TestResult struct {
	IDType   string
	Count    int64
	Duration float64
}

func main() {
	// Create a connection pool
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	// Define the test cases
	tests := []struct {
		functionName string
		idType       string
	}{
		{"test_big_serial_insert", "BigSerial"},
		{"test_serial_uuid_insert", "SerialUUID"},
		{"test_pg_ulid_insert", "PgULID"},
		{"test_pgulid_insert", "PgULIDText"},
		{"test_uuid_varchar_insert", "UUIDVarchar"},
		{"test_uuidv4_insert", "UUIDv4"},
		{"test_uuidv7_insert", "UUIDv7"},
	}

	// Define the row counts to test
	rowCounts := []int64{1000, 10000, 100000}

	// Run the tests and collect results
	var results []TestResult
	for _, test := range tests {
		for _, count := range rowCounts {
			duration, err := runTest(pool, test.functionName, count)
			if err != nil {
				log.Printf("Error running test %s with count %d: %v\n", test.functionName, count, err)
				continue
			}
			results = append(results, TestResult{IDType: test.idType, Count: count, Duration: duration})
		}
	}

	// Print results
	for _, result := range results {
		fmt.Printf("ID Type: %s, Count: %d, Duration: %.2f ms\n", result.IDType, result.Count, result.Duration)
	}
}

// runTest executes a test function in the database and returns the duration in milliseconds
func runTest(pool *pgxpool.Pool, functionName string, count int64) (float64, error) {
	var duration float64
	query := fmt.Sprintf("SELECT %s($1)", functionName)
	start := time.Now()
	err := pool.QueryRow(context.Background(), query, count).Scan(&duration)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start).Milliseconds()
	fmt.Printf("Test %s with count %d completed in %d ms\n", functionName, count, elapsed)
	return float64(elapsed), nil
}
