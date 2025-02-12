package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucsky/cuid"
)

// CUIDGenerator generates CUIDs
type CUIDGenerator struct{}

func NewCUIDGenerator() CUIDGenerator {
	return CUIDGenerator{}
}

func (g CUIDGenerator) Generate() string {
	return cuid.New()
}

func (g CUIDGenerator) CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS cuid_table (id TEXT PRIMARY KEY)")
	return err
}

func (g CUIDGenerator) DropTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS cuid_table")
	return err
}

func (g CUIDGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		_, err := pool.Exec(context.Background(), "INSERT INTO cuid_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g CUIDGenerator) BulkWriteRecords(pool *pgxpool.Pool, count int64) error {
	batch := &pgx.Batch{}
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO cuid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(context.Background(), batch)
	return br.Close()
}

func (g CUIDGenerator) CollectStats(pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(context.Background(), statsQuery, "cuid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
