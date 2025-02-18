package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// NanoIDGenerator generates NanoIDs
type NanoIDGenerator struct{}

var _ IDGenerator = (*NanoIDGenerator)(nil)

func NewNanoIDGenerator() NanoIDGenerator {
	return NanoIDGenerator{}
}

func (g NanoIDGenerator) Generate() string {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return id
}

func (g NanoIDGenerator) Name() string {
	return "NanoID - VARCHAR(21)"
}

func (g NanoIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS nanoid_table (id VARCHAR(21) PRIMARY KEY)")
	return err
}

func (g NanoIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS nanoid_table")
	return err
}

func (g NanoIDGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		_, err := pool.Exec(context.Background(), "INSERT INTO nanoid_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g NanoIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO nanoid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g NanoIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "nanoid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}

func (g NanoIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO nanoid_table (id) VALUES ($1)", g.Generate())
	return err
}
