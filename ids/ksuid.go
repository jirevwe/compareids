package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
)

// KSUIDGenerator generates KSUIDs
type KSUIDGenerator struct{}

func NewKSUIDGenerator() KSUIDGenerator {
	return KSUIDGenerator{}
}

func (g KSUIDGenerator) Generate() string {
	return ksuid.New().String()
}

func (g KSUIDGenerator) CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS ksuid_table (id TEXT PRIMARY KEY)")
	return err
}

func (g KSUIDGenerator) DropTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS ksuid_table")
	return err
}

func (g KSUIDGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		_, err := pool.Exec(context.Background(), "INSERT INTO ksuid_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g KSUIDGenerator) BulkWriteRecords(pool *pgxpool.Pool, count int64) error {
	batch := &pgx.Batch{}
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO ksuid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(context.Background(), batch)
	return br.Close()
}

func (g KSUIDGenerator) CollectStats(pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(context.Background(), statsQuery, "ksuid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
