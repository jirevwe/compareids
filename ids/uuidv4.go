package ids

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UUIDv4Generator generates UUIDv4 IDs
type UUIDv4Generator struct{}

func NewUUIDv4Generator() UUIDv4Generator {
	return UUIDv4Generator{}
}

func (g UUIDv4Generator) Generate() string {
	return uuid.New().String()
}

func (g UUIDv4Generator) CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS uuidv4_table (id TEXT PRIMARY KEY)")
	return err
}

func (g UUIDv4Generator) DropTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS uuidv4_table")
	return err
}

func (g UUIDv4Generator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		_, err := pool.Exec(context.Background(), "INSERT INTO uuidv4_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g UUIDv4Generator) BulkWriteRecords(pool *pgxpool.Pool, count int64) error {
	batch := &pgx.Batch{}
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO uuidv4_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(context.Background(), batch)
	return br.Close()
}

func (g UUIDv4Generator) CollectStats(pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(context.Background(), statsQuery, "uuidv4_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
