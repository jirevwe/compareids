package ids

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UUIDv7GoogleGenerator generates UUIDv7 IDs using the Google UUID package
type UUIDv7GoogleGenerator struct{}

var _ IDGenerator = (*UUIDv7GoogleGenerator)(nil)

func NewUUIDv7GoogleGenerator() UUIDv7GoogleGenerator {
	return UUIDv7GoogleGenerator{}
}

func (g UUIDv7GoogleGenerator) Generate() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (g UUIDv7GoogleGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS uuidv7_google_table (id TEXT PRIMARY KEY)")
	return err
}

func (g UUIDv7GoogleGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS uuidv7_google_table")
	return err
}

func (g UUIDv7GoogleGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO uuidv7_google_table (id) VALUES ($1)", g.Generate())
	return err
}

func (g UUIDv7GoogleGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO uuidv7_google_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g UUIDv7GoogleGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "uuidv7_google_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
