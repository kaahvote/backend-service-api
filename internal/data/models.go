package data

import "database/sql"

type Models struct {
	Sessions SessionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Sessions: SessionModel{DB: db},
	}
}
