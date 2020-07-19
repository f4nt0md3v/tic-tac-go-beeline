package ctx

import "database/sql"

// DbContext contains database pools.
// Add anything related to DB here.
type DbContext struct {
	Postgres *sql.DB
}

func NewDbContext(postgres *sql.DB) *DbContext {
	return &DbContext{Postgres: postgres}
}
