package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// NanoIDGenerator generates NanoIDs
type NanoIDGenerator struct{}

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

func (g NanoIDGenerator) CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS nanoid_table (id TEXT PRIMARY KEY)")
	return err
}

func (g NanoIDGenerator) DropTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS nanoid_table")
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

func (g NanoIDGenerator) BulkWriteRecords(pool *pgxpool.Pool, count int64) error {
	batch := &pgx.Batch{}
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO nanoid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(context.Background(), batch)
	return br.Close()
}

func (g NanoIDGenerator) CollectStats(pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(context.Background(), statsQuery, "nanoid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
