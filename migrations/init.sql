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
RETURNS RECORD AS $$
DECLARE
    rec RECORD;
BEGIN
IF mtype = 'gauge'  then
    insert into gauges values(default, mname, mtype, mvalue)
    on conflict (name, type) do update set value = mvalue
        where gauges.name = mname and gauges.type = mtype;
    select mname, mtype, g.value from gauges g
    where g.name = mname and g.type = mtype into rec;
elsif mtype = 'counter' then 
    insert into counters values(default, mname, mtype, mvalue::bigint)
    on conflict (name, type) do update set delta = (select delta + mvalue::bigint from counters c where c.name = mname and c.type = mtype)
        where counters.name = mname and counters.type = mtype;
    select mname, mtype, delta::double precision from counters c
    where c.name = mname and c.type = mtype into rec;
end if;
return rec;
END
$$ LANGUAGE plpgsql;