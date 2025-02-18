package ids

import (
	"context"
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SnowflakeGenerator generates Snowflake IDs
type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator() SnowflakeGenerator {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to create Snowflake node: %v", err)
	}
	return SnowflakeGenerator{node: node}
}

func (g SnowflakeGenerator) Generate() string {
	return g.node.Generate().String()
}

func (g SnowflakeGenerator) Name() string {
	return "Snowflake - BIGINT"
}

var _ IDGenerator = (*SnowflakeGenerator)(nil)

func (g SnowflakeGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS snowflake_table (id VARCHAR(20) PRIMARY KEY)")
	return err
}

func (g SnowflakeGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS snowflake_table")
	return err
}

func (g SnowflakeGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	// Method removed as it is not part of the IDGenerator interface
	return nil
}

func (g SnowflakeGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.node.Generate().String()
		batch.Queue("INSERT INTO snowflake_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g SnowflakeGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "snowflake_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}

func (g SnowflakeGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO snowflake_table (id) VALUES ($1)", g.node.Generate().String())
	return err
}
