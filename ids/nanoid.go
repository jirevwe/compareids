package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// NanoIDGenerator generates NanoIDs
type NanoIDGenerator struct{}

var _ IDGenerator = (*NanoIDGenerator)(nil)

func NewNanoIDGenerator() *NanoIDGenerator {
	return &NanoIDGenerator{}
}

func (n *NanoIDGenerator) Generate() string {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return id
}

func (n *NanoIDGenerator) Name() string {
	return "NanoID - VARCHAR(21)"
}

func (n *NanoIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS nanoid_table (id VARCHAR(21) PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (n *NanoIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS nanoid_table")
	return err
}

func (n *NanoIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		batch.Queue("INSERT INTO nanoid_table (id, n) VALUES ($1, $2)", n.Generate(), i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (n *NanoIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "nanoid_table", "nanoid_table", "nanoid_table")).Scan(
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

func (n *NanoIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO nanoid_table (id, n) VALUES ($1, $2)", n.Generate(), 1)
	return err
}
