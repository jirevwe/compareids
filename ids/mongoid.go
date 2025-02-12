package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoIDGenerator generates MongoDB ObjectIDs
type MongoIDGenerator struct{}

func NewMongoIDGenerator() MongoIDGenerator {
	return MongoIDGenerator{}
}

func (g MongoIDGenerator) Generate() string {
	return primitive.NewObjectID().Hex()
}

func (g MongoIDGenerator) CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS mongoid_table (id TEXT PRIMARY KEY)")
	return err
}

func (g MongoIDGenerator) DropTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS mongoid_table")
	return err
}

func (g MongoIDGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		_, err := pool.Exec(context.Background(), "INSERT INTO mongoid_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g MongoIDGenerator) BulkWriteRecords(pool *pgxpool.Pool, count int64) error {
	batch := &pgx.Batch{}
	for i := int64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO mongoid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(context.Background(), batch)
	return br.Close()
}

func (g MongoIDGenerator) CollectStats(pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(context.Background(), statsQuery, "mongoid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
