package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
)

// KSUIDGenerator generates KSUIDs
type KSUIDGenerator struct{}

var _ IDGenerator = (*KSUIDGenerator)(nil)

func NewKSUIDGenerator() KSUIDGenerator {
	return KSUIDGenerator{}
}

func (g KSUIDGenerator) Generate() string {
	return ksuid.New().String()
}

func (g KSUIDGenerator) Name() string {
	return "KSUID - VARCHAR(27)"
}

func (g KSUIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS ksuid_table (id VARCHAR(27) PRIMARY KEY)")
	return err
}

func (g KSUIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS ksuid_table")
	return err
}

func (g KSUIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO ksuid_table (id) VALUES ($1)", g.Generate())
	if err != nil {
		return err
	}
	return nil
}

func (g KSUIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		batch.Queue("INSERT INTO ksuid_table (id) VALUES ($1)", g.Generate())
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g KSUIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "ksuid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
