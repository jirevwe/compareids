drop function if exists test_big_serial_insert(count BIGINT);
create or replace function test_big_serial_insert(count BIGINT) returns DOUBLE PRECISION as $$
declare
    v_start double precision;
    v_end   double precision;
begin
    drop table if exists test_big_serial;
    create table if not exists test_big_serial(id bigserial primary key, n bigint not null);
    v_start := extract(epoch from clock_timestamp());
    insert into test_big_serial(n) select g.n from generate_series(1,count) as g(n);
    v_end := extract(epoch from clock_timestamp());
    raise notice 'Time taken to insert %s UUIDv7 records: %s', count, v_end - v_start;
    return v_end - v_start;
end;
$$ language plpgsql;

select
    test_big_serial_insert(100000) as t_100_000,
    test_big_serial_insert(1000000) as t_1_000_000,
    test_big_serial_insert(10000000) as t_10_000_000;
select pg_size_pretty(pg_total_relation_size('test_big_serial')) as total_table_size,
       pg_size_pretty(pg_relation_size('test_big_serial')) as data_size,
       pg_size_pretty(pg_indexes_size('test_big_serial')) as index_size;