package postgres

import "github.com/jackc/pgx/v4"

type postgres struct {
	db *pgx.Conn
}

func (p postgres) DoSomething() error {
	return nil
}