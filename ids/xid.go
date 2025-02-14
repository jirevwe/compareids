package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

// XIDGenerator generates XIDs
type XIDGenerator struct{}

var _ IDGenerator = (*XIDGenerator)(nil)

func NewXIDGenerator() XIDGenerator {
	return XIDGenerator{}
}

func (g XIDGenerator) Generate() string {
	return xid.New().String()
}

func (g XIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS xid_table (id TEXT PRIMARY KEY)")
	return err
}

func (g XIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS xid_table")
	return err
}

func (g XIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO xid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g XIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "xid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}

func (g XIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO xid_table (id) VALUES ($1)", g.Generate())
	return err
}
