// postgres handle read and writes operations to postgres storage
package postgres

import (
	"context"
	"errors"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/metrics"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type readActionByType map[string]func(*postgres, string) (tuples.TupleList, error)
type writeActionByType map[string]func(*postgres, tuples.Tupler) (tuples.TupleList, error)

var readAllMetrics func(p *postgres, mname string) (tuples.TupleList, error) = func(p *postgres, mname string) (tuples.TupleList, error) {
	gauges, _ := readGauges(p, mname)
	counters, nil := readCounters(p, mname)
	gauges = gauges.Merge(counters)
	return gauges, nil
}

var readByType readActionByType = readActionByType{
	"gauge":   readGauges,
	"counter": readCounters,
	"*":       readAllMetrics,
}

var writeByType writeActionByType = writeActionByType{
	"gauge":   writeGauges,
	"counter": writeCounters,
}

// NewPg returnes new instance of *postgres
func NewPg() *postgres {
	return &postgres{
		conn: connection(),
	}
}

type postgres struct {
	conn *pgxpool.Pool
}

// Ping checks for connection with DB
//
// Post-cond: if connection is establish -- nil returned
func (p *postgres) Ping() error {
	if p.conn == nil {
		return errors.New("connection is nil")
	}

	err := p.conn.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// Read reads tuples from database by given query
//
// Pre-cond: given query tuple
// Post-cond: return tuples that satisfies given query
func (p *postgres) Read(state tuples.Tupler) (tuples.TupleList, error) {
	mname, mtype := tuples.ExtractString("name", state), tuples.ExtractString("type", state)

	read := readByType[mtype]
	return read(p, mname)
}

// Write writes given tuple to Database
//
// Pre-cond: given tuple to write
//
// Post-cond: depends on sucsess
// If success then state was written to database and returned written tuple and error = nil
// Otherwise returns nil and error
func (p *postgres) Write(states tuples.TupleList) (tuples.TupleList, error) {
	return p.recWrite(states, tuples.TupleList{})
}

func (p *postgres) assembleCounterState(rows pgx.Rows) (tuples.TupleList, error) {
	defer rows.Close()
	res := tuples.TupleList{}
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

func (p *postgres) assembleGaugeState(rows pgx.Rows) (tuples.TupleList, error) {
	defer rows.Close()
	res := tuples.TupleList{}
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

func connection() *pgxpool.Pool {
	PgURL := server.ServerCfg.DBUrl
	log.Printf("Url:%s", PgURL)
	conn, err := pgxpool.New(context.Background(), PgURL)
	if err != nil {
		log.Printf("conn err :%v", err)
		return nil
	}
	migrate(conn)
	return conn
}

func migrate(c *pgxpool.Pool) {
	// #TODO it must be const
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

func (p *postgres) recWrite(tail tuples.TupleList, acc tuples.TupleList) (tuples.TupleList, error) {
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

func (p *postgres) writeMetric(state tuples.Tupler) (tuples.TupleList, error) {
	mtype := tuples.ExtractString("type", state)
	writeAction := writeByType[mtype]
	return writeAction(p, state)
}

func readGauges(p *postgres, cond string) (tuples.TupleList, error) {
	rows, err := p.conn.Query(context.Background(), READ_METRIC, "gauge", cond)
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

func readCounters(p *postgres, cond string) (tuples.TupleList, error) {
	rows, err := p.conn.Query(context.Background(), READ_METRIC, "counter", cond)
	if err != nil {
		return tuples.TupleList{}, err
	}

	return p.assembleCounterState(rows)
}

func writeGauges(p *postgres, state tuples.Tupler) (tuples.TupleList, error) {
	val := tuples.ExtractFloat64Pointer("value", state)
	if val == nil {
		return tuples.TupleList{}, errors.New("value must exists while writing")
	}

	mname := tuples.ExtractString("name", state)
	mtype := tuples.ExtractString("type", state)
	rows, err := p.conn.Query(context.Background(), WRITE_METRIC, mtype, mname, *val)
	if err != nil {
		return tuples.TupleList{}, err
	}

	return p.assembleGaugeState(rows)
}

func writeCounters(p *postgres, state tuples.Tupler) (tuples.TupleList, error) {
	val := tuples.ExtractInt64Pointer("value", state)
	if val == nil {
		return tuples.TupleList{}, errors.New("value must exists while writing")
	}
	mname := tuples.ExtractString("name", state)
	mtype := tuples.ExtractString("type", state)
	rows, err := p.conn.Query(context.Background(), WRITE_METRIC, mtype, mname, *val)
	if err != nil {
		log.Printf("err counter %v", err)
		return tuples.TupleList{}, err
	}

	return p.assembleCounterState(rows)
}
