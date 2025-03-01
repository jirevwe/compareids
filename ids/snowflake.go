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

func NewSnowflakeGenerator() *SnowflakeGenerator {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to create Snowflake node: %v", err)
	}
	return &SnowflakeGenerator{node: node}
}

func (s *SnowflakeGenerator) Generate() string {
	return s.node.Generate().String()
}

func (s *SnowflakeGenerator) Name() string {
	return "Snowflake - BIGINT"
}

var _ IDGenerator = (*SnowflakeGenerator)(nil)

func (s *SnowflakeGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS snowflake_table (id BIGINT PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (s *SnowflakeGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS snowflake_table")
	return err
}

func (s *SnowflakeGenerator) InsertRecords(pool *pgxpool.Pool, count int64) error {
	// Method removed as it is not part of the IDGenerator interface
	return nil
}

func (s *SnowflakeGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		id := s.node.Generate().Int64()
		batch.Queue("INSERT INTO snowflake_table (id, n) VALUES ($1, $2)", id, i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (s *SnowflakeGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
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

func (s *SnowflakeGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO snowflake_table (id, n) VALUES ($1, $2)", s.node.Generate().Int64(), 1)
	return err
}
