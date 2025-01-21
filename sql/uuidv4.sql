drop function if exists test_uuidv4_insert(count BIGINT);
create or replace function test_uuidv4_insert(count BIGINT) returns DOUBLE PRECISION as $$
declare
    v_start double precision;
    v_end   double precision;
begin
    drop table if exists test_uuidv4;
    create table test_uuidv4(id uuid primary key default gen_random_uuid(), n bigint not null);
    v_start := extract(epoch from clock_timestamp());
    insert into test_uuidv4(n) select g.n from generate_series(1,count) as g(n);
    v_end := extract(epoch from clock_timestamp());
    raise notice 'Time taken to insert %s UUIDv4 records: %s', count, v_end - v_start;
    return v_end - v_start;
end;
$$ language plpgsql;

select test_uuidv4_insert(100000) as t_100_000,
       test_uuidv4_insert(1000000) as t_1_000_000,
       test_uuidv4_insert(10000000) as t_10_000_000;
select pg_size_pretty(pg_total_relation_size('test_uuidv4')) as total_table_size,
       pg_size_pretty(pg_relation_size('test_uuidv4')) as data_size,
       pg_size_pretty(pg_indexes_size('test_uuidv4')) as index_size;