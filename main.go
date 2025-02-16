package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/template"
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
		{"XID", ids.NewXIDGenerator()},
		{"ULID", ids.NewULIDGenerator()},
		{"CUID", ids.NewCUIDGenerator()},
		{"KSUID", ids.NewKSUIDGenerator()},
		{"NanoID", ids.NewNanoIDGenerator()},
		{"UUIDv4", ids.NewUUIDv4Generator()},
		{"UUIDv7", ids.NewUUIDv7Generator()},
		{"TypeID", ids.NewTypeIDGenerator()},
		{"MongoID", ids.NewMongoIDGenerator()},
		{"UUIDv4Db", ids.NewUUIDv4DBGenerator()},
		{"UUIDv7Db", ids.NewUUIDv7DBGenerator()},
		{"Snowflake", ids.NewSnowflakeGenerator()},
		{"BigSerial", ids.NewBigSerialGenerator()},
		{"UUIDv7Google", ids.NewUUIDv7GoogleGenerator()},
		{"ULIDDb", ids.NewULIDDBGenerator()},
		{"ULIDPg", ids.NewULIDPGGenerator()},
	}

	// Define the row counts to test
	rowCounts := []uint64{1000, 10000, 100000, 1000000}

	// Collect all stats
	allStats := make(map[string][]map[string]string)

	// Run the tests and collect results
	var results []TestResult
	for _, test := range tests {
		for _, count := range rowCounts {
			duration, err := runTest(pool, test.generator, count, test.idType, allStats)
			if err != nil {
				log.Printf("Error running test for %s with count %d: %v\n", test.idType, count, err)
				continue
			}
			results = append(results, TestResult{IDType: test.idType, Count: count, Duration: duration})
		}
	}

	// Write all stats to a single JSON file
	file, err := os.Create("all_stats.json")
	if err != nil {
		log.Fatalf("Unable to create JSON file: %v\n", err)
	}
	defer file.Close()

	jsonTemplate := `{
	{{range $idType, $statsList := .}}
	"{{$idType}}": [
		{{range $index, $stats := $statsList}}
		{
			"count": "{{index $stats "count"}}",
			"total_table_size": {{index $stats "total_table_size"}},
			"data_size": {{index $stats "data_size"}},
			"index_size": {{index $stats "index_size"}}
		}{{if not (last $index $statsList)}},{{end}}
		{{end}}
	],
	{{end}}
}`

	// Define a custom function map with a 'last' function
	funcMap := template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == len(a.([]map[string]string))-1
		},
	}

	tmpl, err := template.New("json").Funcs(funcMap).Parse(jsonTemplate)
	if err != nil {
		log.Fatalf("Unable to parse JSON template: %v\n", err)
	}

	fmt.Printf("%+v\n", allStats)
	err = tmpl.Execute(file, allStats)
	if err != nil {
		log.Fatalf("Unable to execute JSON template: %v\n", err)
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
