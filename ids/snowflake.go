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

func (g SnowflakeGenerator) CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS snowflake_table (id TEXT PRIMARY KEY)")
	return err
}

func (g SnowflakeGenerator) DropTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS snowflake_table")
	return err
}

func (g SnowflakeGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	for i := int64(0); i < count; i++ {
		id := g.node.Generate().String()
		_, err := pool.Exec(context.Background(), "INSERT INTO snowflake_table (id) VALUES ($1)", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g SnowflakeGenerator) BulkWriteRecords(pool *pgxpool.Pool, count int64) error {
	batch := &pgx.Batch{}
	for i := int64(0); i < count; i++ {
		id := g.node.Generate().String()
		batch.Queue("INSERT INTO snowflake_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(context.Background(), batch)
	return br.Close()
}

func (g SnowflakeGenerator) CollectStats(pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var count int64
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM snowflake_table").Scan(&count)
	if err != nil {
		return nil, err
	}
	stats["count"] = count
	return stats, nil
}
