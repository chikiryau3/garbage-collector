package pgstorage

import (
	"context"
	"database/sql"
	"fmt"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	_ "github.com/jackc/pgerrcode"
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
	res, err := s.db.QueryContext(ctx, "CREATE TABLE IF NOT EXISTS gauge(name text UNIQUE, value double precision);")
	if err != nil {
		return NewPgError(err)
	}
	err = res.Err()
	if err != nil {
		return NewPgError(err)
	}

	res, err = s.db.QueryContext(ctx, "CREATE TABLE IF NOT EXISTS counter(name text UNIQUE, value bigint);")
	if err != nil {
		return NewPgError(err)
	}
	err = res.Err()
	if err != nil {
		return NewPgError(err)
	}

	return nil
}

func (s *storage) WriteMetric(mtype string, name string, value any) error {
	// todo: think about security issue
	qs := fmt.Sprintf("INSERT INTO %s VALUES ('%s', %v) ON CONFLICT (name) DO UPDATE SET value=%v", mtype, name, value, value)
	res, err := s.db.Query(qs)
	if err != nil {
		return fmt.Errorf("cannot write %s:%v db error %w", name, value, NewPgError(err))
	}
	err = res.Err()
	if err != nil {
		return fmt.Errorf("cannot write %s:%v db error %w", name, value, NewPgError(err))
	}

	return nil
}

func (s *storage) WriteMetrics(mtype string, name string, value any) error {
	// todo: think about security issue
	qs := fmt.Sprintf("INSERT INTO %s VALUES ('%s', %v) ON CONFLICT (name) DO UPDATE SET value=%v", mtype, name, value, value)
	res, err := s.db.Query(qs)
	if err != nil {
		return fmt.Errorf("cannot write %s:%v db error %w", name, value, err)
	}
	err = res.Err()
	if err != nil {
		return fmt.Errorf("cannot write %s:%v db error %w", name, value, err)
	}

	return nil
}

// TODO: return error instead of OK bool

func (s *storage) ReadMetric(mtype string, name string) (any, error) {
	//data, err := s.GetData()
	//fmt.Printf("STORAGE %#v \n", data)
	//fmt.Printf("MTTYPE %s \n", mtype)
	//fmt.Printf("NAME %s \n", name)
	qs := fmt.Sprintf("SELECT * FROM %s WHERE name='%s'", mtype, name)
	//fmt.Printf("QUERY %s \n", qs)

	row := s.db.QueryRow(qs)
	if err := row.Err(); err != nil {
		return nil, NewPgError(err)
	}

	var mName string
	var value any
	err := row.Scan(&mName, &value)
	if err != nil {
		return nil, NewPgError(err)
	}

	return value, nil
}

func (s *storage) GetData() (*metricscollector.StorageData, error) {
	rows, err := s.db.Query("SELECT * FROM gauge, counter")
	if err != nil {
		return nil, NewPgError(err)
	}
	defer rows.Close()

	var data metricscollector.StorageData = map[string]any{}

	for rows.Next() {
		var name string
		var value any
		if err := rows.Scan(&name, &value); err != nil {
			return &data, NewPgError(err)
		}

		data[name] = value
	}

	err = rows.Err()
	if err != nil {
		return nil, NewPgError(err)
	}

	return &data, nil
}

func (s *storage) RunStorageDumper() <-chan error {
	return nil
}

func (s *storage) RestoreFromDump() error {
	return nil
}
