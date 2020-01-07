package database

import (
	"errors"

	"tiomodel/database/postgresql"
)

// GetDBClient 获取指定的数据库连接
// 支持:
// 	postgrs
func GetDBClient(engine, connect string) (TioDb, error) {
	switch engine {
	case "postgres":
		p := &postgresql.TDB_Postgres{}
		if err := p.Init(connect); err != nil {
			return nil, err
		}
		return p, nil
	}

	return nil, errors.New("No Match DB Engine! ")
}
