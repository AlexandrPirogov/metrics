package postgres

const READ_METRIC = "SELECT * from READ_METRIC($1, $2)"

const WRITE_METRIC = "SELECT * FROM WRITE_METRIC($1::varchar(255), $2::varchar(255), $3::double precision)"
