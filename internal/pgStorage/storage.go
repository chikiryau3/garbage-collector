package pgstorage

import (
	"context"
	"database/sql"
	"fmt"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
)

type PgStorage interface {
	metricscollector.Storage

	CheckConnection(ctx context.Context) error
	Init(ctx context.Context) error
}

type Config struct {
}

type storage struct {
	db     *sql.DB
	config *Config
}

func New(db *sql.DB, c *Config) PgStorage {
	return &storage{
		db:     db,
		config: c,
	}
}

func (s *storage) CheckConnection(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *storage) Init(ctx context.Context) error {
	_, err := s.db.QueryContext(ctx, "CREATE TABLE IF NOT EXISTS gauge(name text UNIQUE, value double precision);")
	if err != nil {
		return err
	}

	_, err = s.db.QueryContext(ctx, "CREATE TABLE IF NOT EXISTS counter(name text UNIQUE, value bigint);")
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) WriteMetric(mtype string, name string, value any) error {
	// todo: think about security issue
	//qs := fmt.Sprintf("UPDATE %s SET value=%v WHERE name=%s", mtype, value, name)
	qs := fmt.Sprintf("INSERT INTO %s VALUES ('%s', %v) ON CONFLICT (name) DO UPDATE SET value=%v", mtype, name, value, value)
	//qs := fmt.Sprintf("INSERT INTO %s (name, value) VALUES ('%s', %v)", mtype, name, value)
	//fmt.Printf("\nQUERY STRING %s \n", qs)
	//_, err := s.db.Query(qs)
	_, err := s.db.Query(qs)
	if err != nil {
		return fmt.Errorf("cannot write %s:%v db error %w", name, value, err)
	}

	return nil
}

// TODO: return error instead of OK bool

func (s *storage) ReadMetric(mtype string, name string) (any, bool) {
	row := s.db.QueryRow("SELECT * FROM $1 WHERE name=$2", mtype, name)
	if err := row.Err(); err != nil {
		return nil, false
	}

	var value any
	err := row.Scan(&value)
	if err != nil {
		return nil, false
	}

	return value, true
}

func (s *storage) GetData() (*metricscollector.StorageData, error) {
	rows, err := s.db.Query("SELECT * FROM gauge, counter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data metricscollector.StorageData

	for rows.Next() {
		var name string
		var value any
		if err := rows.Scan(&name, &value); err != nil {
			return &data, err
		}

		data[name] = value
	}

	return &data, nil
}

func (s *storage) RunStorageDumper() <-chan error {
	return nil
}

func (s *storage) RestoreFromDump() error {
	return nil
}
