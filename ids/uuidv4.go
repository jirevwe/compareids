package ids

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UUIDv4Generator generates UUIDv4 IDs
type UUIDv4Generator struct{}

var _ IDGenerator = (*UUIDv4Generator)(nil)

func NewUUIDv4Generator() *UUIDv4Generator {
	return &UUIDv4Generator{}
}

func (u *UUIDv4Generator) Generate() string {
	return uuid.NewString()
}

func (u *UUIDv4Generator) Name() string {
	return "UUIDv4 - UUID"
}

func (u *UUIDv4Generator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS uuidv4_table (id UUID PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (u *UUIDv4Generator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS uuidv4_table")
	return err
}

func (u *UUIDv4Generator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		batch.Queue("INSERT INTO uuidv4_table (id, n) VALUES ($1, $2)", u.Generate(), i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (u *UUIDv4Generator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "uuidv4_table", "uuidv4_table", "uuidv4_table")).Scan(
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

func (u *UUIDv4Generator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO uuidv4_table (id, n) VALUES ($1, $2)", u.Generate(), 1)
	return err
}
