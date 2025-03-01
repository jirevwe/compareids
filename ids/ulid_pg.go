package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ULIDPgGenerator generates ULIDs using the pg-ulid extension
type ULIDPgGenerator struct{}

var _ IDGenerator = (*ULIDPgGenerator)(nil)

func NewULIDPGGenerator() *ULIDPgGenerator {
	return &ULIDPgGenerator{}
}

func (u *ULIDPgGenerator) Generate() string {
	// ULID is generated by the database, so this might return an empty string or a placeholder.
	return ""
}

func (u *ULIDPgGenerator) Name() string {
	return "ULID (PG) - ULID"
}

func (u *ULIDPgGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	err := u.LoadULIDFunction(ctx, pool)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS ulid_pg_table (id ulid PRIMARY KEY DEFAULT gen_ulid(), n BIGINT NOT NULL)")
	return err
}

func (u *ULIDPgGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS ulid_pg_table")
	return err
}

func (u *ULIDPgGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO ulid_pg_table (n) VALUES (1)")
	return err
}

func (u *ULIDPgGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	_, err := pool.Exec(ctx, "INSERT INTO ulid_pg_table (n) SELECT g.n FROM generate_series(1, $1) AS g(n)", count)
	return err
}

func (u *ULIDPgGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "ulid_pg_table", "ulid_pg_table", "ulid_pg_table")).Scan(
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

func (u *ULIDPgGenerator) LoadULIDFunction(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	CREATE EXTENSION IF NOT EXISTS ulid with schema public;
	`
	_, err := pool.Exec(ctx, sql)
	return err
}
