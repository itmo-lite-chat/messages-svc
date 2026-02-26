package example_storage

import "github.com/doug-martin/goqu/v9"

type Storage struct {
	db *goqu.Database
}

func NewStorage(db *goqu.Database) *Storage {
	return &Storage{db: db}
}