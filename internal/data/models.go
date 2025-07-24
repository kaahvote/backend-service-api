package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Sessions SessionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Sessions: SessionModel{DB: db},
	}
}
