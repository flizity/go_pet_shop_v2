package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Storage struct {
	conn clickhouse.Conn
}

func New() (*Storage, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return nil, err
	}
	return &Storage{conn: conn}, nil
}

func (s *Storage) Close() error {
	return s.conn.Close()
}
