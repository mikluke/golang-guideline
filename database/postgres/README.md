# Postgres

# Index

* [Drivers](#drivers)
* [Orm or not](#orm-or-not)
* [Migrations](#migrations)
* [Example](#example)

# Topics


## Drivers
Right now there are only two production ready and highly used solutions:
* https://github.com/jackc/pgx
* https://github.com/lib/pq

Nevertheless `lib/pq` usage is decreasing and even maintainers recommend to switch to `pgx`.

## Orm or not
Overall no. There is no production ready and stable ORMs for go and in long term perspective using one of them will make your code hardly supportable. For simplifying work with sql we recommend to use https://github.com/jmoiron/sqlx.

## Migrations
[What is migrations](https://en.wikipedia.org/wiki/Schema_migration)
* https://github.com/rubenv/sql-migrate
* https://github.com/golang-migrate/migrate

We also recommend not to mix migration and application logic. In other words don't apply migrations from your go code. Do it inside your ci/cd or even manually.

## Example
For our examples we are using `pgx` driver with sqlx extensions.

### Client setup
```go
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	conn sqlx.ExtContext
}

type Config struct {
	URL             string
	MaxConns        int          
	MaxIdleConns    int          
	MaxConnIdleTime time.Duration
}

func New(ctx context.Context, cfg Config) (*Postgres, error) {
	db, err := postgresDB(ctx, cfg.URL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	return &Postgres{conn: db}, nil
}

func postgresDB(ctx context.Context, url string) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres url: %w", err)
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}
```

### Fetch data
```go
// It's a good practice to have a storage model separated from the domain one.
type user struct {
    id   int    `db:"id"`
    type string `db:"type"`
    name string `db:"name"`
}

func (p *Postgres) Users(ctx context.Context, limit, offset int) ([]domain.User, error) {
	rows, err := p.conn.QueryxContext(
		ctx,
		`SELECT id, type, name FROM users LIMIT $1 OFFSET $2`,
    	limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query users: %w", err)
	}

	var users []domain.User
	for rows.Next() {
		var u user
		if err := rows.StructScan(&u); err != nil {
			return nil, fmt.Errorf("could not scan user: %w", err)
		}
		users = append(users, mapToDomainUser(&u))
	}

	return users, nil
}
```

