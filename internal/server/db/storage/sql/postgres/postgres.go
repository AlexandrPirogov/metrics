package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
	"os"

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
func (p *Postgres) Write(states tuples.TupleList) (tuples.TupleList, error) {
	return p.recWrite(states, tuples.TupleList{})
}

func (p *Postgres) recWrite(tail tuples.TupleList, acc tuples.TupleList) (tuples.TupleList, error) {
	if !tail.Next() {
		return acc, nil
	}

	head, tail := tail.HeadTail()
	written, err := p.writeMetric(head)
	if err != nil {
		return tuples.TupleList{}, err
	}
	return p.recWrite(tail, acc.Merge(written))
}

func (p *Postgres) writeMetric(state tuples.Tupler) (tuples.TupleList, error) {
	written := tuples.TupleList{}
	mtype := tuples.ExtractString("type", state)
	switch mtype {
	case "gauge":
		res, err := p.WriteGauges(state)
		if err != nil {
			return tuples.TupleList{}, err
		}
		written = written.Merge(res)
	case "counter":
		res, err := p.WriteCounters(state)
		if err != nil {
			return tuples.TupleList{}, err
		}
		written = written.Merge(res)
	}
	return written, nil
}

// Read reads tuples from database by given query
//
// Pre-cond: given query tuple
// Post-cond: return tuples that satisfies given query
func (p *Postgres) Read(state tuples.Tupler) (tuples.TupleList, error) {
	mname := tuples.ExtractString("name", state)
	mtype := tuples.ExtractString("type", state)

	switch mtype {
	case "gauge":
		res, _ := p.ReadGauges(mname)
		return res, nil
	case "counter":
		res, _ := p.ReadCounters(mname)
		return res, nil
	case "*":
		gauges, _ := p.ReadGauges(mname)
		counters, nil := p.ReadCounters(mname)
		gauges = gauges.Merge(counters)
		return gauges, nil
	default:
		return tuples.TupleList{}, nil
	}
}

func Ping() error {
	conn := connection()
	if conn == nil {
		return errors.New("connection is nil")
	}
	defer conn.Close(context.Background())
	err := conn.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) ReadGauges(cond string) (tuples.TupleList, error) {
	query := fmt.Sprintf("SELECT * from READ_METRIC('gauge', '%s')", cond)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return tuples.TupleList{}, err
	}

	defer rows.Close()

	var res = tuples.TupleList{}
	for rows.Next() {
		var toScan metrics.GaugeState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return tuples.TupleList{}, err
		}
		res = res.Add(toScan)
	}
	return res, nil
}

func (p *Postgres) ReadCounters(cond string) (tuples.TupleList, error) {
	query := fmt.Sprintf("SELECT * FROM READ_METRIC('counter', '%s')", cond)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return tuples.TupleList{}, err
	}
	defer rows.Close()

	var res = tuples.TupleList{}

	for rows.Next() {
		var toScan metrics.CounterState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return tuples.TupleList{}, err
		}
		res = res.Add(toScan)
	}
	return res, nil
}

func (p *Postgres) ReadMetrics() (tuples.TupleList, error) {
	query := "SELECT * FROM READ_METRICS()"
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return tuples.TupleList{}, err
	}

	defer rows.Close()

	var res = tuples.TupleList{}

	for rows.Next() {
		var toScan metrics.CounterState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return tuples.TupleList{}, err
		}
		res = res.Add(toScan)
	}
	return res, nil
}

func (p *Postgres) WriteGauges(state tuples.Tupler) (tuples.TupleList, error) {
	val := tuples.ExtractFloat64Pointer("value", state)
	if val == nil {
		return tuples.TupleList{}, errors.New("value must exists while writing")
	}

	mname := tuples.ExtractString("name", state)
	mtype := tuples.ExtractString("type", state)
	rows, err := p.conn.Query(context.Background(),
		"SELECT * FROM WRITE_METRIC($1::varchar(255), $2::varchar(255), $3::double precision)", mtype, mname, *val)
	if err != nil {
		return tuples.TupleList{}, err
	}

	defer rows.Close()

	var res = tuples.TupleList{}
	for rows.Next() {
		var toScan metrics.GaugeState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			return tuples.TupleList{}, err
		}
		res = res.Add(toScan)
	}
	return res, nil
}

func (p *Postgres) WriteCounters(state tuples.Tupler) (tuples.TupleList, error) {

	val := tuples.ExtractInt64Pointer("value", state)
	if val == nil {
		return tuples.TupleList{}, errors.New("value must exists while writing")
	}

	mname := tuples.ExtractString("name", state)
	mtype := tuples.ExtractString("type", state)
	log.Printf("writing counter\n")
	rows, err := p.conn.Query(context.Background(),
		"SELECT * FROM WRITE_METRIC($1::varchar(255), $2::varchar(255), $3::double precision)", mtype, mname, *val)
	if err != nil {
		log.Printf("err counter %v", err)
		return tuples.TupleList{}, err
	}

	defer rows.Close()

	var res = tuples.TupleList{}

	for rows.Next() {
		var toScan metrics.CounterState
		err := rows.Scan(&toScan.Name, &toScan.Type, &toScan.Value)
		if err != nil {
			log.Printf("error scan %v", err)
			return tuples.TupleList{}, err
		}
		res = res.Add(toScan)
	}
	return res, nil
}

func connection() *pgx.Conn {
	PgURL := server.ServerCfg.DBUrl
	log.Printf("Url:%s", PgURL)
	conn, err := pgx.Connect(context.Background(), PgURL)
	if err != nil {
		log.Printf("conn err :%v", err)
		return nil
	}
	Migrate(conn)
	return conn
}

func Migrate(c *pgx.Conn) {
	path := "internal/server/db/storage/sql/postgres/init.sql"
	body, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error while read%v", err)
	}

	sql := string(body)
	_, err = c.Exec(context.Background(), sql)
	if err != nil {
		log.Printf("Error while migrate%v", err)
	}
}
