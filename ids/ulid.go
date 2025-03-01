package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

// ULIDGenerator generates ULID IDs
type ULIDGenerator struct{}

var _ IDGenerator = (*ULIDGenerator)(nil)

func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

func (u *ULIDGenerator) Generate() string {
	return ulid.Make().String()
}

func (u *ULIDGenerator) Name() string {
	return "ULID - VARCHAR(26)"
}

func (u *ULIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS ulid_table (id TEXT PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (u *ULIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS ulid_table")
	return err
}

func (u *ULIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		batch.Queue("INSERT INTO ulid_table (id, n) VALUES ($1, $2)", u.Generate(), i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (u *ULIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "ulid_table", "ulid_table", "ulid_table")).Scan(
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

func (u *ULIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO ulid_table (id, n) VALUES ($1, $2)", u.Generate(), 1)
	return err
}
