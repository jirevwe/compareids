package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	Count    uint64
	Duration float64
}

// TestCase represents a single ID generator test case
type TestCase struct {
	generator ids.IDGenerator
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
	tests := []TestCase{
		{ids.NewXIDGenerator()},
		{ids.NewULIDGenerator()},
		{ids.NewCUIDGenerator()},
		{ids.NewKSUIDGenerator()},
		{ids.NewNanoIDGenerator()},
		{ids.NewUUIDv4Generator()},
		{ids.NewUUIDv7Generator()},
		{ids.NewTypeIDGenerator()},
		{ids.NewMongoIDGenerator()},
		{ids.NewUUIDv4DBGenerator()},
		{ids.NewUUIDv7DBGenerator()},
		{ids.NewSnowflakeGenerator()},
		{ids.NewBigSerialGenerator()},
		{ids.NewUUIDv7GoogleGenerator()},
		{ids.NewULIDDBGenerator()},
		{ids.NewULIDPGGenerator()},
	}

	// Define the row counts to test
	rowCounts := []uint64{
		1_000,
		10_000,
		100_000,
		// 1_000_000,
		// 10_000_000,
	}

	// Collect all stats
	dbSizeStats := make(map[string][]map[string]string)
	insertDurationStats := make(map[string][]float64)

	// Run the tests and collect stats
	var results []TestResult
	for _, test := range tests {
		for _, count := range rowCounts {
			duration, err := runTest(pool, test.generator, count, test.generator.Name(), dbSizeStats)
			if err != nil {
				log.Printf("Error running test for %s with count %d: %v\n", test.generator.Name(), count, err)
				continue
			}

			insertDurationStats[test.generator.Name()] = append(insertDurationStats[test.generator.Name()], duration)
			results = append(results, TestResult{IDType: test.generator.Name(), Count: count, Duration: duration})

			log.Printf("Took %vms to insert %d records for %s\n", duration, count, test.generator.Name())
		}

		// Drop the table after all row counts have been tested for this generator
		if err := test.generator.DropTable(context.Background(), pool); err != nil {
			log.Printf("Error dropping table for %s: %v\n", test.generator.Name(), err)
		}
	}

	// Create a map for the template data
	templateData := struct {
		Data    map[string][]map[string]interface{}
		IDTypes []string
	}{
		Data:    make(map[string][]map[string]interface{}),
		IDTypes: make([]string, 0),
	}

	// Process the results for the template
	for _, test := range tests {
		idType := test.generator.Name()
		templateData.IDTypes = append(templateData.IDTypes, idType)
		templateData.Data[idType] = make([]map[string]interface{}, 0)

		// Get all results for this ID type
		for _, result := range results {
			if result.IDType == idType {
				// Get the corresponding size stats
				var sizeStats map[string]string
				for _, stats := range dbSizeStats[idType] {
					if stats["count"] == fmt.Sprintf("%d", result.Count) {
						sizeStats = stats
						break
					}
				}

				// Combine duration and size stats
				dataPoint := map[string]interface{}{
					"count":            result.Count,
					"duration":         result.Duration,
					"total_table_size": sizeStats["total_table_size"],
					"data_size":        sizeStats["data_size"],
					"index_size":       sizeStats["index_size"],
				}
				templateData.Data[idType] = append(templateData.Data[idType], dataPoint)
			}
		}
	}

	// Write template data to JSON file
	templateDataFile, err := os.Create("template_data.json")
	if err != nil {
		log.Fatalf("Unable to create template data JSON file: %v\n", err)
	}
	defer templateDataFile.Close()

	encoder := json.NewEncoder(templateDataFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(templateData); err != nil {
		log.Fatalf("Unable to write template data to JSON file: %v\n", err)
	}

	// Write all stats to a single JSON file
	file, err := os.Create("db_size_stats.json")
	if err != nil {
		log.Fatalf("Unable to create JSON file: %v\n", err)
	}
	defer file.Close()

	// Write the stats directly as JSON
	encoder = json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dbSizeStats); err != nil {
		log.Fatalf("Unable to write stats to JSON file: %v\n", err)
	}
}

// runTest generates IDs on the client side and inserts them into the database
func runTest(pool *pgxpool.Pool, g ids.IDGenerator, count uint64, idType string, allStats map[string][]map[string]string) (float64, error) {
	start := time.Now()
	ctx := context.Background()

	// Begin a transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	err = g.DropTable(ctx, pool)
	if err != nil {
		return 0, err
	}

	// create the table
	err = g.CreateTable(ctx, pool)
	if err != nil {
		return 0, err
	}

	// err = g.InsertRecord(ctx, pool)
	// if err != nil {
	// 	return 0, err
	// }

	// Insert generated IDs
	err = g.BulkWriteRecords(ctx, pool, count)
	if err != nil {
		return 0, err
	}

	// Collect stats after inserting records
	stats, err := g.CollectStats(ctx, pool)
	if err != nil {
		return 0, err
	}

	// Convert stats to map[string]string
	convertedStats := make(map[string]string)
	for k, v := range stats {
		if str, ok := v.(string); ok {
			convertedStats[k] = str
		}
	}

	// Add the count to the convertedStats map
	convertedStats["count"] = fmt.Sprintf("%d", count)

	// Collect all stats
	allStats[idType] = append(allStats[idType], convertedStats)

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	elapsed := time.Since(start).Milliseconds()
	return float64(elapsed), nil
}
