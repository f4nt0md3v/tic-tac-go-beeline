package context

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/websocketx"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/repositories"
)

// DbContext contains database pools.
// Add anything related to DB here.
type DbContext struct {
	Postgres *sql.DB
}

func NewDbContext(postgres *sql.DB) *DbContext {
	return &DbContext{Postgres: postgres}
}

// AppContext contains local context:
// database context or anything else handlers need to access.
type AppContext struct {
	DbContext *DbContext
	GameRepo  *repositories.GameRepo
	Logger    *zap.SugaredLogger
	Pool      *websocketx.Pool
}

func NewAppContext(d *DbContext, g *repositories.GameRepo, l *zap.SugaredLogger, p *websocketx.Pool) *AppContext {
	return &AppContext{
		DbContext: d,
		GameRepo:  g,
		Logger:    l,
		Pool:      p,
	}
}
