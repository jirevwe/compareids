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
}

const statsQuery = `SELECT
	pg_size_pretty(pg_total_relation_size($1)) AS total_table_size, 
	pg_size_pretty(pg_relation_size($1)) AS data_size, 
	pg_size_pretty(pg_indexes_size($1)) AS index_size;`
