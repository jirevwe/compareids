
drop function if exists test_serial_uuid_insert(count BIGINT);
create or replace function test_serial_uuid_insert(count BIGINT) returns DOUBLE PRECISION as $$
declare
    v_start double precision;
    v_end   double precision;
begin
    drop table if exists test_serial_uuid;
    create table test_serial_uuid(id bigserial primary key, u uuid not null default gen_random_uuid());
    v_start := extract(epoch from clock_timestamp());
    insert into test_serial_uuid(id) select g.id from generate_series(1,count) as g(id);
    v_end := extract(epoch from clock_timestamp());
    raise notice 'Time taken to insert %s BigSerial x UUID records: %s', count, v_end - v_start;
    return v_end - v_start;
end;
$$ language plpgsql;

select
    test_serial_uuid_insert(100000) as t_100_000,
    test_serial_uuid_insert(1000000) as t_1_000_000,
    test_serial_uuid_insert(10000000) as t_10_000_000;

select pg_size_pretty(pg_total_relation_size('test_serial_uuid')) as total_table_size,
       pg_size_pretty(pg_relation_size('test_serial_uuid')) as data_size,
       pg_size_pretty(pg_indexes_size('test_serial_uuid')) as index_size;
