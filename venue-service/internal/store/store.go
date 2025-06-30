package store

import "database/sql"

type Store struct {
}

func NewStore(postgres *sql.DB) Store {
	return Store{}
}
