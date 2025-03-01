package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucsky/cuid"
)

// CUIDGenerator generates CUIDs
type CUIDGenerator struct{}

var _ IDGenerator = (*CUIDGenerator)(nil)

func NewCUIDGenerator() *CUIDGenerator {
	return &CUIDGenerator{}
}

func (c *CUIDGenerator) Name() string {
	return "CUID - VARCHAR(25)"
}

func (c *CUIDGenerator) Generate() string {
	return cuid.New()
}

func (c *CUIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS cuid_table (id VARCHAR(25) PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (c *CUIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS cuid_table")
	return err
}

func (c *CUIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		batch.Queue("INSERT INTO cuid_table (id, n) VALUES ($1, $2)", c.Generate(), i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (c *CUIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "cuid_table", "cuid_table", "cuid_table")).Scan(
		&tableStats.TotalTableSize,
		&tableStats.DataSize,
		&tableStats.IndexSize,
		&tableStats.InternalPages,
		&tableStats.LeafPages,
		&tableStats.Density,
		&tableStats.Fragmentation,
	)
	if err != nil {
		return nil, err
	}

	stats["total_table_size"] = tableStats.TotalTableSize
	stats["data_size"] = tableStats.DataSize
	stats["index_size"] = tableStats.IndexSize
	stats["index_internal_pages"] = tableStats.InternalPages
	stats["index_leaf_pages"] = tableStats.LeafPages
	stats["index_density"] = tableStats.Density
	stats["index_fragmentation"] = tableStats.Fragmentation

	// Calculate the ratio of internal pages to leaf pages
	stats["index_internal_to_leaf_ratio"] = float64(tableStats.InternalPages) / float64(tableStats.LeafPages)

	return stats, nil
}

func (c *CUIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO cuid_table (id, n) VALUES ($1, $2)", c.Generate(), 1)
	return err
}
