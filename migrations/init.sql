CREATE TABLE COUNTERS(
    ID SERIAL,
    NAME VARCHAR(255),
    TYPE VARCHAR(255),
    DELTA BIGINT,
    PRIMARY KEY(NAME),
    UNIQUE(NAME, TYPE)
);

CREATE TABLE GAUGES(
    ID SERIAL,
    NAME VARCHAR(255),
    TYPE VARCHAR(255),
    VALUE DOUBLE PRECISION,
    PRIMARY KEY(NAME),
    UNIQUE(NAME, TYPE)
);

CREATE OR REPLACE FUNCTION WRITE_METRIC(mtype VARCHAR(255), mname VARCHAR(255), mvalue DOUBLE PRECISION) 
RETURNS TABLE (
    rtype VARCHAR(255),
    rname VARCHAR(255),
    rvalue DOUBLE PRECISION
)  AS $$
BEGIN
IF mtype = 'gauge'  then
    insert into gauges values(default, mname, mtype, mvalue)
    on conflict (name, type) do update set value = mvalue
        where gauges.name = mname and gauges.type = mtype;
    return query select mname as rname, mtype as rtype, g.value as rvalue from gauges g
    where g.name = mname and g.type = mtype ;
elsif mtype = 'counter' then 
    insert into counters values(default, mname, mtype, mvalue::bigint)
    on conflict (name, type) do update set delta = (select delta + mvalue::bigint from counters c where c.name = mname and c.type = mtype)
        where counters.name = mname and counters.type = mtype;
     return query  select mname as rname, mtype as rtype, delta::double precision as value from counters c
    where c.name = mname and c.type = mtype ;
end if;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION READ_METRIC(mtype VARCHAR(255), mname VARCHAR(255)) 
RETURNS TABLE (
    rtype VARCHAR(255),
    rname VARCHAR(255),
    rvalue DOUBLE PRECISION
)  AS $$
BEGIN
IF mtype = 'gauge'  then
    return query select mname as rname, mtype as rtype, g.value as rvalue from gauges g
    where g.name = mname and g.type = mtype ;
elsif mtype = 'counter' then
    return query select mname as rname, mtype as rtype, delta::double precision as value from counters c
    where c.name = mname and c.type = mtype ;
end if;

END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION READ_METRICS(mtype VARCHAR(255)) 
RETURNS TABLE (
    rtype VARCHAR(255),
    rname VARCHAR(255),
    rvalue DOUBLE PRECISION
)  AS $$
BEGIN
IF mtype = 'gauge'  then
    return query select mname as rname, mtype as rtype, g.value as rvalue from gauges g
    where g.type = mtype ;
elsif mtype = 'counter' then
    return query  select mname as rname, mtype as rtype, delta::double precision as value from counters c
    where c.type = mtype ;
end if;

END
$$ LANGUAGE plpgsql;