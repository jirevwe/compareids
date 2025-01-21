/*
 * MIT License
 *
 * Copyright (c) 2023-2024 Fabio Lima
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

/**
 * Returns a time-ordered UUID with Unix Epoch (UUIDv7).
 *
 * Reference: https://www.rfc-editor.org/rfc/rfc9562.html
 * Reference: https://gist.github.com/fabiolimace/515a0440e3e40efeb234e12644a6a346
 *
 * MIT License.
 *
 */
create or replace function uuid7() returns uuid as $$
declare
begin
    return uuid7(clock_timestamp());
end $$ language plpgsql;

create or replace function uuid7(p_timestamp timestamp with time zone) returns uuid as $$
declare

    v_time double precision := null;

    v_unix_t bigint := null;
    v_rand_a bigint := null;
    v_rand_b bigint := null;

    v_unix_t_hex varchar := null;
    v_rand_a_hex varchar := null;
    v_rand_b_hex varchar := null;

    c_milli double precision := 10^3;  -- 1 000
    c_micro double precision := 10^6;  -- 1 000 000
    c_scale double precision := 4.096; -- 4.0 * (1024 / 1000)

    c_version bigint := x'0000000000007000'::bigint; -- RFC-9562 version: b'0111...'
    c_variant bigint := x'8000000000000000'::bigint; -- RFC-9562 variant: b'10xx...'

begin

    v_time := extract(epoch from p_timestamp);

    v_unix_t := trunc(v_time * c_milli);
    v_rand_a := trunc((v_time * c_micro - v_unix_t * c_milli) * c_scale);
--     v_rand_b := secure_random_bigint(); -- use when pgcrypto extension is installed
    v_rand_b := trunc(random() * 2^30)::bigint << 32 | trunc(random() * 2^32)::bigint;

    v_unix_t_hex := lpad(to_hex(v_unix_t), 12, '0');
    v_rand_a_hex := lpad(to_hex((v_rand_a | c_version)::bigint), 4, '0');
    v_rand_b_hex := lpad(to_hex((v_rand_b | c_variant)::bigint), 16, '0');

    return (v_unix_t_hex || v_rand_a_hex || v_rand_b_hex)::uuid;

end $$ language plpgsql;

-- select uuid7() uuid, clock_timestamp()-statement_timestamp() time_taken;

-- DO $$
-- DECLARE
-- v_start double precision;
-- v_end   double precision;
-- BEGIN
-- 	v_start := extract(epoch from clock_timestamp());
--  insert into exp_uuidv7(n) select g.n from generate_series(1,1000000) as g(n);
-- 	v_end := extract(epoch from clock_timestamp());
-- 	raise notice 'Time taken: %s', v_end - v_start;
-- END
-- $$ language plpgsql;

drop function  if exists  test_uuidv7_insert(count BIGINT);
create or replace function test_uuidv7_insert(count BIGINT) returns DOUBLE PRECISION as $$
declare
    v_start double precision;
    v_end   double precision;
begin
    drop table if exists exp_uuidv7;
    create table exp_uuidv7(id uuid primary key default uuid7(), n bigint not null);
    v_start := extract(epoch from clock_timestamp());
    insert into exp_uuidv7(n) select g.n from generate_series(1,count) as g(n);
    v_end := extract(epoch from clock_timestamp());
    raise notice 'Time taken to insert %s UUIDv7 records: %s', count, v_end - v_start;
    return v_end - v_start;
end;
$$ language plpgsql;

select test_uuidv7_insert(100000) as t_100_000,
       test_uuidv7_insert(1000000) as t_1_000_000,
       test_uuidv7_insert(10000000) as t_10_000_000;

select pg_size_pretty(pg_total_relation_size('exp_uuidv7')) as total_table_size,
       pg_size_pretty(pg_relation_size('exp_uuidv7')) as data_size,
       pg_size_pretty(pg_indexes_size('exp_uuidv7')) as index_size;