package ids

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UUIDv7Generator generates UUIDv7 IDs
type UUIDv7Generator struct{}

var _ IDGenerator = (*UUIDv7Generator)(nil)

func NewUUIDv7Generator() UUIDv7Generator {
	return UUIDv7Generator{}
}

func (g UUIDv7Generator) Generate() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (g UUIDv7Generator) Name() string {
	return "UUIDv7 - UUID"
}

func (g UUIDv7Generator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS uuidv7_table (id UUID PRIMARY KEY)")
	return err
}

func (g UUIDv7Generator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS uuidv7_table")
	return err
}

func (g UUIDv7Generator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	// Method removed as it is not part of the IDGenerator interface
	return nil
}

func (g UUIDv7Generator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO uuidv7_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g UUIDv7Generator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "uuidv7_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}

func (g UUIDv7Generator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO uuidv7_table (id) VALUES ($1)", g.Generate())
	return err
}
