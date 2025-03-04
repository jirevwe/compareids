package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jirevwe/compareids/ids"
)

// GetIDGenerator returns the ID generator for the given ID type
func GetIDGenerator(idType string) (ids.IDGenerator, error) {
	switch idType {
	case "bigserial":
		return ids.NewBigSerialGenerator(), nil
	case "snowflake":
		return ids.NewSnowflakeGenerator(), nil
	case "uuidv4":
		return ids.NewUUIDv4Generator(), nil
	case "uuidv4-db":
		return ids.NewUUIDv4DBGenerator(), nil
	case "uuidv7":
		return ids.NewUUIDv7Generator(), nil
	case "uuidv7-db":
		return ids.NewUUIDv7DBGenerator(), nil
	case "uuidv7-google":
		return ids.NewUUIDv7GoogleGenerator(), nil
	case "ulid":
		return ids.NewULIDGenerator(), nil
	case "ulid-db":
		return ids.NewULIDDBGenerator(), nil
	case "ulid-pg":
		return ids.NewULIDPGGenerator(), nil
	case "xid":
		return ids.NewXIDGenerator(), nil
	case "cuid":
		return ids.NewCUIDGenerator(), nil
	case "ksuid":
		return ids.NewKSUIDGenerator(), nil
	case "nanoid":
		return ids.NewNanoIDGenerator(), nil
	case "typeid":
		return ids.NewTypeIDGenerator(), nil
	case "mongoid":
		return ids.NewMongoIDGenerator(), nil
	default:
		return nil, fmt.Errorf("unknown ID type: %s", idType)
	}
}

// GetAllIDTypes returns a list of all supported ID types
func GetAllIDTypes() []string {
	return []string{
		"bigserial",
		"snowflake",
		"uuidv4",
		"uuidv4-db",
		"uuidv7",
		"uuidv7-db",
		"uuidv7-google",
		"ulid",
		"ulid-db",
		"ulid-pg",
		"xid",
		"cuid",
		"ksuid",
		"nanoid",
		"typeid",
		"mongoid",
	}
}

// GetDefaultRowCounts returns the default row counts to test
func GetDefaultRowCounts() []uint64 {
	return []uint64{
		1_000,
		10_000,
		100_000,
		1_000_000,
	}
}

// RunTest generates IDs and inserts them into the database
func RunTest(pool *pgxpool.Pool, g ids.IDGenerator, count uint64) (float64, map[string]string, error) {
	start := time.Now()
	ctx := context.Background()

	// Begin a transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback(ctx)

	err = g.DropTable(ctx, pool)
	if err != nil {
		return 0, nil, err
	}

	// create the table
	err = g.CreateTable(ctx, pool)
	if err != nil {
		return 0, nil, err
	}

	// Measure system resources during the bulk write operation
	systemMetrics, err := MeasureSystemResources(func() error {
		return g.BulkWriteRecords(ctx, pool, count)
	})
	if err != nil {
		return 0, nil, err
	}

	// Collect stats after inserting records
	stats, err := g.CollectStats(ctx, pool)
	if err != nil {
		return 0, nil, err
	}

	// Convert stats to map[string]string
	convertedStats := make(map[string]string)
	for k, v := range stats {
		if str, ok := v.(string); ok {
			convertedStats[k] = str
		}

		if str, ok := v.(int); ok {
			convertedStats[k] = fmt.Sprintf("%d", str)
		}

		if str, ok := v.(int64); ok {
			convertedStats[k] = fmt.Sprintf("%d", str)
		}

		if str, ok := v.(float64); ok {
			convertedStats[k] = fmt.Sprintf("%f", str)
		}
	}

	// Add the count to the convertedStats map
	convertedStats["count"] = fmt.Sprintf("%d", count)

	// Add system metrics to the stats
	systemMetricsMap := systemMetrics.AsMap()
	for k, v := range systemMetricsMap {
		convertedStats[k] = v
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, nil, err
	}

	duration := float64(time.Since(start).Milliseconds())
	return duration, convertedStats, nil
}

// SaveTestResult saves the test result to a JSON file
func SaveTestResult(result TestResult) error {
	// Create the results directory if it doesn't exist
	if err := os.MkdirAll(ResultsDir, 0755); err != nil {
		return err
	}

	// Create the file path
	filePath := filepath.Join(ResultsDir, fmt.Sprintf("%s_%d.json", result.IDType, result.Count))

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the result to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// LoadTestResults loads all test results from the results directory
func LoadTestResults() ([]TestResult, error) {
	// Create the results directory if it doesn't exist
	if err := os.MkdirAll(ResultsDir, 0755); err != nil {
		return nil, err
	}

	// Get all JSON files in the results directory
	files, err := filepath.Glob(filepath.Join(ResultsDir, "*.json"))
	if err != nil {
		return nil, err
	}

	// Load each file
	var results []TestResult
	for _, file := range files {
		// Open the file
		f, err := os.Open(file)
		if err != nil {
			log.Printf("Error opening file %s: %v", file, err)
			continue
		}

		// Decode the JSON
		var result TestResult
		if err := json.NewDecoder(f).Decode(&result); err != nil {
			f.Close()
			log.Printf("Error decoding file %s: %v", file, err)
			continue
		}
		f.Close()

		results = append(results, result)
	}

	return results, nil
}
