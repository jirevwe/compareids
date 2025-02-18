package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucsky/cuid"
)

// CUIDGenerator generates CUIDs
type CUIDGenerator struct{}

var _ IDGenerator = (*CUIDGenerator)(nil)

func NewCUIDGenerator() CUIDGenerator {
	return CUIDGenerator{}
}

func (g CUIDGenerator) Name() string {
	return "CUID - VARCHAR(25)"
}

func (g CUIDGenerator) Generate() string {
	return cuid.New()
}

func (g CUIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS cuid_table (id VARCHAR(25) PRIMARY KEY)")
	return err
}

func (g CUIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS cuid_table")
	return err
}

func (g CUIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO cuid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g CUIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "cuid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}

func (g CUIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO cuid_table (id) VALUES ($1)", g.Generate())
	return err
}
