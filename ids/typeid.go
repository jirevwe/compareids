package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.jetify.com/typeid"
)

// TypeIDGenerator generates TypeIDs
type TypeIDGenerator struct{}

var _ IDGenerator = (*TypeIDGenerator)(nil)

func NewTypeIDGenerator() *TypeIDGenerator {
	return &TypeIDGenerator{}
}

type CustomPrefix struct{}

func (CustomPrefix) Prefix() string {
	return ""
}

func (t *TypeIDGenerator) Generate() string {
	id, err := typeid.New[typeid.TypeID[CustomPrefix]]()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (t *TypeIDGenerator) Name() string {
	return "TypeID - VARCHAR(27)"
}

func (t *TypeIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS typeid_table (id VARCHAR(27) PRIMARY KEY, n BIGINT NOT NULL)")
	return err
}

func (t *TypeIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS typeid_table")
	return err
}

func (t *TypeIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO typeid_table (id, n) VALUES ($1, $2)", t.Generate(), 1)
	return err
}

func (t *TypeIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(1); i <= count; i++ {
		id := t.Generate()
		batch.Queue("INSERT INTO typeid_table (id, n) VALUES ($1, $2)", id, i)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (t *TypeIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "typeid_table", "typeid_table", "typeid_table")).Scan(
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
