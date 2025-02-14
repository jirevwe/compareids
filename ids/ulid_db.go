package ids

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ULIDDBGenerator generates ULIDs using the database
type ULIDDBGenerator struct{}

var _ IDGenerator = (*ULIDDBGenerator)(nil)

func NewULIDDBGenerator() ULIDDBGenerator {
	return ULIDDBGenerator{}
}

func (g ULIDDBGenerator) Generate() string {
	// ULID is generated by the database, so this might return an empty string or a placeholder.
	return ""
}

func (g ULIDDBGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	err := g.LoadULIDFunction(ctx, pool)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS ulid_table (id TEXT PRIMARY KEY DEFAULT generate_ulid(), n BIGINT NOT NULL)")
	return err
}

func (g ULIDDBGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS ulid_table")
	return err
}

func (g ULIDDBGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO ulid_table (n) VALUES (1)")
	return err
}

func (g ULIDDBGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	_, err := pool.Exec(ctx, "INSERT INTO ulid_table (n) SELECT g.n FROM generate_series(1, $1) AS g(n)", count)
	return err
}

func (g ULIDDBGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)
	var totalTableSize, dataSize, indexSize string
	err := pool.QueryRow(ctx, statsQuery, "ulid_table").Scan(&totalTableSize, &dataSize, &indexSize)
	if err != nil {
		return nil, err
	}
	stats["total_table_size"] = totalTableSize
	stats["data_size"] = dataSize
	stats["index_size"] = indexSize
	return stats, nil
}

func (g ULIDDBGenerator) LoadULIDFunction(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	CREATE EXTENSION IF NOT EXISTS pgcrypto with schema public;

	CREATE OR REPLACE FUNCTION generate_ulid() RETURNS TEXT AS $$
	DECLARE
		encoding BYTEA = '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
		timestamp BYTEA = E'\\000\\000\\000\\000\\000\\000';
		output TEXT = '';
		unix_time BIGINT;
		ulid BYTEA;
	BEGIN
		unix_time = (EXTRACT(EPOCH FROM CLOCK_TIMESTAMP()) * 1000)::BIGINT;
		timestamp = SET_BYTE(timestamp, 0, (unix_time >> 40)::BIT(8)::INTEGER);
		timestamp = SET_BYTE(timestamp, 1, (unix_time >> 32)::BIT(8)::INTEGER);
		timestamp = SET_BYTE(timestamp, 2, (unix_time >> 24)::BIT(8)::INTEGER);
		timestamp = SET_BYTE(timestamp, 3, (unix_time >> 16)::BIT(8)::INTEGER);
		timestamp = SET_BYTE(timestamp, 4, (unix_time >> 8)::BIT(8)::INTEGER);
		timestamp = SET_BYTE(timestamp, 5, unix_time::BIT(8)::INTEGER);
		ulid = timestamp || gen_random_bytes(10);
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 0) & 224) >> 5));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 0) & 31)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 1) & 248) >> 3));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 1) & 7) << 2) | ((GET_BYTE(ulid, 2) & 192) >> 6)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 2) & 62) >> 1));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 2) & 1) << 4) | ((GET_BYTE(ulid, 3) & 240) >> 4)));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 3) & 15) << 1) | ((GET_BYTE(ulid, 4) & 128) >> 7)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 4) & 124) >> 2));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 4) & 3) << 3) | ((GET_BYTE(ulid, 5) & 224) >> 5)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 5) & 31)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 6) & 248) >> 3));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 6) & 7) << 2) | ((GET_BYTE(ulid, 7) & 192) >> 6)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 7) & 62) >> 1));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 7) & 1) << 4) | ((GET_BYTE(ulid, 8) & 240) >> 4)));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 8) & 15) << 1) | ((GET_BYTE(ulid, 9) & 128) >> 7)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 9) & 124) >> 2));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 9) & 3) << 3) | ((GET_BYTE(ulid, 10) & 224) >> 5)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 10) & 31)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 11) & 248) >> 3));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 11) & 7) << 2) | ((GET_BYTE(ulid, 12) & 192) >> 6)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 12) & 62) >> 1));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 12) & 1) << 4) | ((GET_BYTE(ulid, 13) & 240) >> 4)));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 13) & 15) << 1) | ((GET_BYTE(ulid, 14) & 128) >> 7)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 14) & 124) >> 2));
		output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 14) & 3) << 3) | ((GET_BYTE(ulid, 15) & 224) >> 5)));
		output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 15) & 31)));
		RETURN output;
	END $$ LANGUAGE plpgsql VOLATILE;
	`
	_, err := pool.Exec(ctx, sql)
	return err
}
