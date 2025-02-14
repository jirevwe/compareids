package ids

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.jetify.com/typeid"
)

// TypeIDGenerator generates TypeIDs
type TypeIDGenerator struct{}

var _ IDGenerator = (*TypeIDGenerator)(nil)

func NewTypeIDGenerator() TypeIDGenerator {
	return TypeIDGenerator{}
}

type CustomPrefix struct{}

func (CustomPrefix) Prefix() string {
	return "custom"
}

func (g TypeIDGenerator) Generate() string {
	id, err := typeid.New[typeid.TypeID[CustomPrefix]]()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (g TypeIDGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS typeid_table (id TEXT PRIMARY KEY)")
	return err
}

func (g TypeIDGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS typeid_table")
	return err
}

func (g TypeIDGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO typeid_table (id) VALUES ($1)", g.Generate())
	return err
}

func (g TypeIDGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	batch := &pgx.Batch{}
	for i := uint64(0); i < count; i++ {
		id := g.Generate()
		batch.Queue("INSERT INTO typeid_table (id) VALUES ($1)", id)
	}
	br := pool.SendBatch(ctx, batch)
	return br.Close()
}

func (g TypeIDGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "typeid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}
