package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jirevwe/compareids/ids"
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
		idType    string
		generator ids.IDGenerator
	}{
		{"UUIDv4", ids.NewUUIDv4Generator()},
		{"UUIDv7", ids.NewUUIDv7Generator()},
		{"ULID", ids.NewULIDGenerator()},
		{"CUID", ids.NewCUIDGenerator()},
		{"KSUID", ids.NewKSUIDGenerator()},
		{"NanoID", ids.NewNanoIDGenerator()},
		{"XID", ids.NewXIDGenerator()},
		{"Snowflake", ids.NewSnowflakeGenerator()},
		{"MongoID", ids.NewMongoIDGenerator()},
	}

	// Define the row counts to test
	rowCounts := []int64{1000, 10000, 100000}

	// Run the tests and collect results
	var results []TestResult
	for _, test := range tests {
		for _, count := range rowCounts {
			duration, err := runTest(pool, test.generator, count)
			if err != nil {
				log.Printf("Error running test for %s with count %d: %v\n", test.idType, count, err)
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

// runTest generates IDs on the client side and inserts them into the database
func runTest(pool *pgxpool.Pool, generator ids.IDGenerator, count int64) (float64, error) {
	start := time.Now()

	// Begin a transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	// Prepare the insert statement
	stmt, err := tx.Prepare(context.Background(), "insert_stmt", "INSERT INTO public.test_table (id, n) VALUES ($1, $2)")
	if err != nil {
		return 0, err
	}

	// Insert generated IDs
	for i := int64(1); i <= count; i++ {
		id := generator.Generate()
		_, err := tx.Exec(context.Background(), stmt.Name, id, i)
		if err != nil {
			return 0, err
		}
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	elapsed := time.Since(start).Milliseconds()
	fmt.Printf("Test completed with count %d in %d ms\n", count, elapsed)
	return float64(elapsed), nil
}
