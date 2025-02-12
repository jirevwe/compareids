package ids

import "github.com/jackc/pgx/v5/pgxpool"

// IDGenerator is an interface for generating IDs
type IDGenerator interface {
	Generate() string
	CreateTable(pool *pgxpool.Pool) error
	DropTable(pool *pgxpool.Pool) error
	InsertRecords(pool *pgxpool.Pool, count int64) error
	BulkWriteRecords(pool *pgxpool.Pool, count int64) error
	CollectStats(pool *pgxpool.Pool) (map[string]any, error)
}

const statsQuery = `SELECT
	pg_size_pretty(pg_total_relation_size($1)) AS total_table_size, 
	pg_size_pretty(pg_relation_size($1)) AS data_size, 
	pg_size_pretty(pg_indexes_size($1)) AS index_size;`
