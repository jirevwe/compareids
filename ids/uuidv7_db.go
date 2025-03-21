package ids

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UUIDv7DBGenerator generates UUIDv7 IDs using the database
type UUIDv7DBGenerator struct{}

var _ IDGenerator = (*UUIDv7DBGenerator)(nil)

func NewUUIDv7DBGenerator() *UUIDv7DBGenerator {
	return &UUIDv7DBGenerator{}
}

func (u *UUIDv7DBGenerator) Generate() string {
	// UUIDv7 is generated by the database
	return ""
}

func (u *UUIDv7DBGenerator) Name() string {
	return "UUIDv7 (DB) - UUID"
}

func (u *UUIDv7DBGenerator) CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	err := u.LoadUUID7Function(ctx, pool)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS uuidv7_db_table (id UUID PRIMARY KEY DEFAULT uuid7(), n BIGINT NOT NULL)")
	return err
}

func (u *UUIDv7DBGenerator) DropTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS uuidv7_db_table")
	return err
}

func (u *UUIDv7DBGenerator) InsertRecord(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "INSERT INTO uuidv7_db_table (n) VALUES (1)")
	return err
}

func (u *UUIDv7DBGenerator) BulkWriteRecords(ctx context.Context, pool *pgxpool.Pool, count uint64) error {
	_, err := pool.Exec(ctx, "INSERT INTO uuidv7_db_table (n) SELECT g.n FROM generate_series(1, $1) AS g(n)", count)
	return err
}

func (u *UUIDv7DBGenerator) CollectStats(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	stats := make(map[string]any)

	err := LoadPGStatTuple(ctx, pool)
	if err != nil {
		return nil, err
	}

	var tableStats TableStats

	err = pool.QueryRow(ctx, fmt.Sprintf(fmtStatsQuery, "uuidv7_db_table", "uuidv7_db_table", "uuidv7_db_table")).Scan(
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

// LoadUUID7Function Create a PL/PgSQL function that returns a time-ordered UUID with Unix Epoch (UUIDv7).
// Reference: https://gist.github.com/fabiolimace/515a0440e3e40efeb234e12644a6a346
// RFC: https://www.rfc-editor.org/rfc/rfc9562.html
func (u *UUIDv7DBGenerator) LoadUUID7Function(ctx context.Context, pool *pgxpool.Pool) error {
	sql := `
	CREATE OR REPLACE FUNCTION uuid7() RETURNS uuid AS $$
	DECLARE
	BEGIN
		RETURN uuid7(clock_timestamp());
	END $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION uuid7(p_timestamp timestamp with time zone) RETURNS uuid AS $$
	DECLARE
		v_time DOUBLE PRECISION := NULL;
		v_unix_t BIGINT := NULL;
		v_rand_a BIGINT := NULL;
		v_rand_b BIGINT := NULL;
		v_unix_t_hex VARCHAR := NULL;
		v_rand_a_hex VARCHAR := NULL;
		v_rand_b_hex VARCHAR := NULL;
		c_milli DOUBLE PRECISION := 10^3;
		c_micro DOUBLE PRECISION := 10^6;
		c_scale DOUBLE PRECISION := 4.096;
		c_version BIGINT := x'0000000000007000'::BIGINT;
		c_variant BIGINT := x'8000000000000000'::BIGINT;
	BEGIN
		v_time := extract(epoch FROM p_timestamp);
		v_unix_t := trunc(v_time * c_milli);
		v_rand_a := trunc((v_time * c_micro - v_unix_t * c_milli) * c_scale);
		v_rand_b := trunc(random() * 2^30)::BIGINT << 32 | trunc(random() * 2^32)::BIGINT;
		v_unix_t_hex := lpad(to_hex(v_unix_t), 12, '0');
		v_rand_a_hex := lpad(to_hex((v_rand_a | c_version)::BIGINT), 4, '0');
		v_rand_b_hex := lpad(to_hex((v_rand_b | c_variant)::BIGINT), 16, '0');
		RETURN (v_unix_t_hex || v_rand_a_hex || v_rand_b_hex)::uuid;
	END $$ LANGUAGE plpgsql;
	`
	_, err := pool.Exec(ctx, sql)
	return err
}
