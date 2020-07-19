package repositories

import (
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models/game"
)

const (
	queryCreateNewGame string = `
		INSERT INTO games(game_id, first_user_id, state) VALUES ($1, $2, $3);
	`
	queryGetGameById string = `
		SELECT
		       game_id,
		       first_user_id,
		       second_user_id,
		       state,
		       last_move_user_id,
		       created_at,
		       last_modified_at
		FROM games WHERE game_id = $1;
	`
	queryUpdateGameWithId string = `
		UPDATE games
		SET
		    first_user_id = $2,
		    second_user_id = $3,
		    state = $4,
		    last_move_user_id = $5,
		    last_modified_at = now()
		WHERE game_id = $1;
	`
)

type GameRepository interface {
	Create(string, string) (*game.Game, error)
	FindByGameID(ID string) (*game.Game, error)
	Update(game *game.Game) error
}

type GameRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewGameRepo(db *sql.DB, log *zap.SugaredLogger) *GameRepo {
	return &GameRepo{db: db, logger: log}
}

func (g *GameRepo) Create(gameId, userId string) (*game.Game, error) {
	initState := "2,2,2,2,2,2,2,2,2"

	stmt, err := g.db.Prepare(queryCreateNewGame)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			g.logger.Error(err)
		}
	}()

	_, err = stmt.Exec(gameId, userId, initState)
	if err != nil {
		return nil, err
	}

	return &game.Game{
		GameId:         gameId,
		FirstUserId:    userId,
		SecondUserId:   "",
		State:          initState,
		LastMoveUserId: "",
	}, nil
}

type GameRaw struct {
	ID             int
	GameId         string
	FirstUserId    sql.NullString
	SecondUserId   sql.NullString
	State          string
	LastMoveUserId sql.NullString
	CreatedAt      *time.Time
	LastModifiedAt *time.Time
}

func (g *GameRaw) toGame() *game.Game {
	return &game.Game{
		ID:             g.ID,
		GameId:         g.GameId,
		FirstUserId:    g.FirstUserId.String,
		SecondUserId:   g.SecondUserId.String,
		State:          g.State,
		LastMoveUserId: g.LastMoveUserId.String,
		CreatedAt:      g.CreatedAt,
		LastModifiedAt: g.LastModifiedAt,
	}
}

func (g *GameRepo) FindByGameID(id string) (*game.Game, error) {
	stmt, err := g.db.Prepare(queryGetGameById)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			g.logger.Error(err)
		}
	}()

	var gameRaw GameRaw
	row := stmt.QueryRow(id)
	if err := row.Scan(
		&gameRaw.GameId,
		&gameRaw.FirstUserId,
		&gameRaw.SecondUserId,
		&gameRaw.State,
		&gameRaw.LastMoveUserId,
		&gameRaw.CreatedAt,
		&gameRaw.LastModifiedAt,
	); err != nil {
		g.logger.Errorf("error when trying to get user by id: %s", err)
		return nil, err
	}
	gm := gameRaw.toGame()
	return gm, nil
}

func (g *GameRepo) Update(game *game.Game) error {
	stmt, err := g.db.Prepare(queryUpdateGameWithId)
	if err != nil {
		return err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			g.logger.Error(err)
		}
	}()

	r, err := stmt.Exec(
		game.GameId,
		game.FirstUserId,
		game.SecondUserId,
		game.State,
		game.LastMoveUserId,
	)
	if err != nil {
		return err
	}

	n, err := r.RowsAffected()
	if err != nil || n == 0 {
		g.logger.Errorf("error while trying to update game: %s", err)
	}
	return nil
}
