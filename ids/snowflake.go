package ids

import (
	"context"
	"fmt"
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

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "snowflake_table", "snowflake_table", "snowflake_table")).Scan(
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

func (g SnowflakeGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO snowflake_table (id) VALUES ($1)", g.node.Generate().String())
	return err
}
