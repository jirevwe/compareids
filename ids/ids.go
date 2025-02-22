package ids

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// IDGenerator is an interface for generating IDs
type IDGenerator interface {
	Generate() string
	CreateTable(ctx context.Context, pool *pgxpool.Pool) error
	DropTable(ctx context.Context, pool *pgxpool.Pool) error
	InsertRecord(ctx context.Context, pool *pgxpool.Pool) error
	BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, recordsWritten uint64) error
	CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error)
	Name() string
}

// TableStats holds all the statistics for a table and its index
type TableStats struct {
	TotalTableSize string  `json:"total_table_size" db:"total_table_size"`
	DataSize       string  `json:"data_size" db:"data_size"`
	IndexSize      int64   `json:"index_size" db:"index_size"`
	InternalPages  int64   `json:"index_internal_pages" db:"index_internal_pages"`
	LeafPages      int64   `json:"index_leaf_pages" db:"index_leaf_pages"`
	Density        float64 `json:"index_density" db:"index_density"`
	Fragmentation  float64 `json:"index_fragmentation" db:"index_fragmentation"`
}

// LoadPGStatTuple ensures the pgstattuple extension is installed
func LoadPGStatTuple(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pgstattuple")
	return err
}

const fmtStatsQuery = `WITH table_stats AS (
	SELECT
		pg_total_relation_size('%s')::text AS total_table_size, 
		pg_relation_size('%s')::text AS data_size
),
index_stats AS (
	SELECT
		index_size::bigint,
		internal_pages::bigint,
		leaf_pages::bigint,
		avg_leaf_density::decimal,
		leaf_fragmentation::decimal
	FROM pgstatindex((
		SELECT i.indexrelid::regclass::text
		FROM pg_index i
		JOIN pg_class c ON c.oid = i.indrelid
		WHERE c.relname = '%s'
		AND i.indisprimary
		LIMIT 1
	))
)
SELECT 
	total_table_size,
	data_size,
	COALESCE(index_size, 0) as index_size,
	COALESCE(internal_pages, 0) as index_internal_pages,
	COALESCE(leaf_pages, 0) as index_leaf_pages,
	COALESCE(avg_leaf_density, 0.0) as index_density,
	COALESCE(leaf_fragmentation, 0.0) as index_fragmentation
FROM table_stats t
LEFT JOIN index_stats i ON true;`
