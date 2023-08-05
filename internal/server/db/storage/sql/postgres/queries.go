package postgres

const ReadMetric = "SELECT * from READ_METRIC($1, $2)"

const WriteMetric = "SELECT * FROM WRITE_METRIC($1::varchar(255), $2::varchar(255), $3::double precision)"
