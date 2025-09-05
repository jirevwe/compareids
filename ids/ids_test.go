package ids

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestStatsCollection(t *testing.T) {
	ctx := context.Background()

	connString := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	generator := NewUUIDv4Generator()

	// Clean up and recreate table
	_ = generator.DropTable(ctx, pool)
	err = generator.CreateTable(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}

	// Insert test data
	err = generator.BulkWriteRecords(ctx, pool, 1000)
	if err != nil {
		t.Fatal(err)
	}

	// Collect stats
	stats, err := generator.CollectStats(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}

	// Verify stats
	assert.NotEmpty(t, stats["total_table_size"], "total_table_size should not be empty")
	assert.NotEmpty(t, stats["data_size"], "data_size should not be empty")
	assert.Greater(t, stats["index_size"].(int64), int64(0), "index_size should be greater than 0")
	assert.Greater(t, stats["index_internal_pages"].(int64), int64(0), "index_internal_pages should be greater than 0")
	assert.Greater(t, stats["index_leaf_pages"].(int64), int64(0), "index_leaf_pages should be greater than 0")
	assert.Greater(t, stats["index_density"].(float64), float64(0), "index_density should be greater than 0")
	assert.Greater(t, stats["index_fragmentation"].(float64), float64(0), "index_fragmentation should be greater than 0")

	t.Logf("Stats collected: %+v", stats)
}
