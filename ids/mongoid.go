package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoIDGenerator generates MongoDB ObjectIDs
type MongoIDGenerator struct{}

var _ IDGenerator = (*MongoIDGenerator)(nil)

func NewMongoIDGenerator() *MongoIDGenerator {
	return &MongoIDGenerator{}
}

func (m *MongoIDGenerator) Generate() string {
	return primitive.NewObjectID().Hex()
}

func (m *MongoIDGenerator) Name() string {
	return "MongoDB ObjectID - VARCHAR(24)"
}

func (m *MongoIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS mongoid_table (id VARCHAR(24) PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (m *MongoIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS mongoid_table")
	return err
}

func (m *MongoIDGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := m.Generate()
		_, err := pool.Exec(context.Background(), "INSERT INTO mongoid_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		id := m.Generate()
		batch.Queue("INSERT INTO mongoid_table (id, n) VALUES ($1, $2)", id, i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (m *MongoIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "mongoid_table", "mongoid_table", "mongoid_table")).Scan(
		&tableStats.TotalTableSize,
		&tableStats.DataSize,
		&tableStats.IndexSize,
		&tableStats.InternalPages,
		&tableStats.LeafPages,
		&tableStats.Density,
		&tableStats.Fragmentation,
	)
	if err != nil {
		return nil, err
	}

	stats["total_table_size"] = tableStats.TotalTableSize
	stats["data_size"] = tableStats.DataSize
	stats["index_size"] = tableStats.IndexSize
	stats["index_internal_pages"] = tableStats.InternalPages
	stats["index_leaf_pages"] = tableStats.LeafPages
	stats["index_density"] = tableStats.Density
	stats["index_fragmentation"] = tableStats.Fragmentation

	// Calculate the ratio of internal pages to leaf pages
	stats["index_internal_to_leaf_ratio"] = float64(tableStats.InternalPages) / float64(tableStats.LeafPages)

	return stats, nil
}

func (m *MongoIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO mongoid_table (id, n) VALUES ($1, $2)", m.Generate(), 1)
	return err
}
