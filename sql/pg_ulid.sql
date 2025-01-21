-- https://github.com/andrielfn/pg-ulid
CREATE EXTENSION IF NOT EXISTS ulid with schema public;

drop function if exists test_pg_ulid_insert(count BIGINT);
create or replace function test_pg_ulid_insert(count BIGINT) returns DOUBLE PRECISION as $$
declare
    v_start double precision;
    v_end   double precision;
begin
    drop table if exists test_pg_ulid;
    create table if not exists test_pg_ulid(id ulid primary key default gen_ulid(), n bigint not null);
    v_start := extract(epoch from clock_timestamp());
    insert into test_pg_ulid(n) select g.n from generate_series(1,count) as g(n);
    v_end := extract(epoch from clock_timestamp());
    raise notice 'Time taken to insert %s ULID records: %s', count, v_end - v_start;
    return v_end - v_start;
end;
$$ language plpgsql;

select
    test_pg_ulid_insert(100000)   as t_100_000,
    test_pg_ulid_insert(1000000)  as t_1_000_000,
    test_pg_ulid_insert(10000000) as t_10_000_000;

select pg_size_pretty(pg_total_relation_size('test_pg_ulid')) as total_table_size,
       pg_size_pretty(pg_relation_size('test_pg_ulid')) as data_size,
       pg_size_pretty(pg_indexes_size('test_pg_ulid')) as index_size;
