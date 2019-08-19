package datasource

import (
	"errors"
	"mvc/datamodels"
)

type Engine uint32

const (
	Memory Engine = iota
	Bolt
	MySql
)

func loadUsers(engine Engine) (map[int64]datamodels.User, error) {
	if engine != Memory {
		return nil, errors.New("请使用map类型作为数据源")
	}
	return make(map[int64]datamodels.User), nil
}
