package ctx

import (
	"go.uber.org/zap"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/repositories"
)

// AppContext contains local context:
// database context or anything else handlers need to access.
type AppContext struct {
	DbContext *DbContext
	GameRepo  *repositories.GameRepo
	Logger    *zap.SugaredLogger
	Pool      *Pool
}

func NewAppContext(d *DbContext, g *repositories.GameRepo, l *zap.SugaredLogger, p *Pool) *AppContext {
	return &AppContext{
		DbContext: d,
		GameRepo:  g,
		Logger:    l,
		Pool:      p,
	}
}
