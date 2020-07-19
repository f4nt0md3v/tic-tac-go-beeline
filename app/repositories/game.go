package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models"
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

type GameRepo struct {
	db *sql.DB
}
type GameRepository interface {
	Create(string, string) (*models.Game, error)
	FindByGameID(ID string) (*models.Game, error)
	Update(game *models.Game) error
}

func NewGameRepo(db *sql.DB) *GameRepo {
	return &GameRepo{db: db}
}

func (g *GameRepo) Create(gameId, userId string) (*models.Game, error) {
	initState := "2,2,2,2,2,2,2,2,2"

	stmt, err := g.db.Prepare(queryCreateNewGame)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	r, err := stmt.Exec(gameId, userId, initState)
	if err != nil {
		return nil, err
	}

	lastId, err := r.LastInsertId()
	if err != nil {
		fmt.Printf("error while trying to get last inserted id: %s", err)
	}

	return &models.Game{
		ID:             int(lastId),
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

func (g *GameRaw) toGame() *models.Game {
	return &models.Game{
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

func (g *GameRepo) FindByGameID(id string) (*models.Game, error) {
	stmt, err := g.db.Prepare(queryGetGameById)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println(err)
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
		fmt.Printf("error when trying to get user by id: %s", err)
		return nil, err
	}
	game := gameRaw.toGame()
	return game, nil
}

func (g *GameRepo) Update(game *models.Game) error {
	stmt, err := g.db.Prepare(queryUpdateGameWithId)
	if err != nil {
		return err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println(err)
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
		fmt.Printf("error while trying to update game: %s", err)
	}
	return nil
}
