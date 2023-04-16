package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"

	"github.com/jackc/pgx/v5"
)

func NewPg() *Postgres {
	return &Postgres{
		conn: connection(),
	}
}

type Postgres struct {
	conn *pgx.Conn
}

// Write writes given tuple to Database
//
// Pre-cond: given tuple to write
//
// Post-cond: depends on sucsess
// If success then state was written to database and returned written tuple and error = nil
// Otherwise returns nil and error
func (p *Postgres) Write(state tuples.Tupler) (tuples.Tupler, error) {

	mname := tuples.ExtractString("name", state)
	mtype := tuples.ExtractString("type", state)

	switch mtype {
	case "gauge":
		val := tuples.ExtractFloat64Pointer("value", state)
		res, err := p.WriteGauges(mname, val)
		if err != nil {
			return nil, err
		}
		return res[0], nil
	case "counter":
		val := tuples.ExtractInt64Pointer("value", state)
		res, err := p.WriteCounteres(mname, val)
		if err != nil {
			return nil, err
		}
		return res[0], nil
	default:
		return nil, errors.New("given type not exists")
	}
}

// Read reads tuples from database by given query
//
// Pre-cond: given query tuple
// Post-cond: return tuples that satisfies given query
func (p *Postgres) Read(state tuples.Tupler) ([]tuples.Tupler, error) {
	mname := tuples.ExtractString("type", state)
	mtype := tuples.ExtractString("name", state)

	switch mtype {
	case "gauge":
		res, _ := p.ReadGauges(mname)
		return res, nil
	case "counter":
		res, _ := p.ReadCounteres(mname)
		return res, nil
	case "*":
		gauges, _ := p.ReadGauges(mname)
		counters, nil := p.ReadCounteres(mname)
		gauges = append(gauges, counters...)
		return gauges, nil
	default:
		return nil, errors.New("bad")
	}
}

func Ping() error {
	conn := connection()

	defer conn.Close(context.Background())
	err := conn.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) ReadGauges(cond string) ([]tuples.Tupler, error) {
	query := fmt.Sprintf("SELECT READ_METRIC('gauge', '%s')", cond)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res = []tuples.Tupler{}
	for rows.Next() {
		var toScan metrics.GaugeState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, toScan)
	}
	return res, nil
}

func (p *Postgres) ReadCounteres(cond string) ([]tuples.Tupler, error) {
	query := fmt.Sprintf("SELECT * FROM READ_METRIC('counter', '%s')", cond)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res = []tuples.Tupler{}

	for rows.Next() {
		var toScan metrics.CounterState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, toScan)
	}
	return res, nil
}

func (p *Postgres) ReadMetrics() ([]tuples.Tupler, error) {
	query := "SELECT * FROM READ_METRICS()"
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res = []tuples.Tupler{}

	for rows.Next() {
		var toScan metrics.CounterState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, toScan)
	}
	return res, nil
}

func (p *Postgres) WriteGauges(cond string, val *float64) ([]tuples.Tupler, error) {
	if val == nil {
		return nil, errors.New("value must exists while writing")
	}

	rows, err := p.conn.Query(context.Background(), "SELECT * FROM WRITE_METRIC('gauge'::varchar(255), $1::varchar(255), $2::double precision)", cond, *val)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res = []tuples.Tupler{}
	for rows.Next() {
		var toScan metrics.GaugeState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, toScan)
	}
	return res, nil
}

func (p *Postgres) WriteCounteres(cond string, val *int64) ([]tuples.Tupler, error) {
	if val == nil {
		return nil, errors.New("value must exists while writing")
	}

	rows, err := p.conn.Query(context.Background(), "SELECT * FROM WRITE_METRIC('counter'::varchar(255), $1::varchar(255), $2::double precision)", cond, *val)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res = []tuples.Tupler{}

	for rows.Next() {
		var toScan metrics.CounterState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, toScan)
	}
	return res, nil
}

func connection() *pgx.Conn {
	PgURL := server.ServerCfg.DBUrl
	log.Printf("Url:%s", PgURL)
	conn, err := pgx.Connect(context.Background(), PgURL)
	if err != nil {
		log.Fatalf("%v", err)
		return nil
	}
	return conn
}
