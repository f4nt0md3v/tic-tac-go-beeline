package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/context"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/websocketx"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/repositories"
)

var (
	upg = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// CheckOrigin: func(r *http.Request) bool { return true }, // TODO: comment out or remove next line on production
	}
)

const (
	CmdGenerateNewGame   = "GENERATE_NEW_GAME"
	CmdJoinGame          = "JOIN_GAME"
	CmdNewMove           = "NEW_MOVE"
	CmdOpponentConnected = "OPPONENT_CONNECTED"
)

// serveWs handles websocket requests from the peer.
func ServeWs(c *context.AppContext, w http.ResponseWriter, r *http.Request) {
	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &websocketx.Client{Pool: c.Pool, Conn: conn, Send: make(chan []byte, 256)}
	client.Pool.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func WebsocketHandler(c *context.AppContext, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.Logger.Error(err)
		}
	}()

	c.Logger.Info("Client connected...")
	client := &websocketx.Client{
		Pool: c.Pool,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Pool.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func OldWebsocketHandler(c echo.Context) error {
	// Upgrade HTTP connection to WebSocket
	ws, err := upg.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := ws.Close(); err != nil {
			c.Logger().Error(err)
		}
	}()

	c.Logger().Print("Client connected...")

	for {
		var req models.Request
		// Receive a message using the codec
		if parseErr := ws.ReadJSON(&req); err != nil {
			c.Logger().Error(parseErr)
			if err := ws.WriteJSON(parseErr); err != nil {
				c.Logger().Error(err)
			}
			return err
		}
		c.Logger().Printf("Received message: %s", req)

		switch req.Command {
		case CmdGenerateNewGame:
			c.Logger().Print("Generating a new game...")
			gameInfo, err := GenerateNewGame(c)
			if err != nil {
				errResp := models.Response{
					Code:  http.StatusInternalServerError,
					Error: err.Error(),
				}
				if err := ws.WriteJSON(errResp); err != nil {
					c.Logger().Error(err)
				}
				return err
			}

			resp := models.Response{
				Code:     http.StatusCreated,
				Command:  CmdGenerateNewGame,
				GameInfo: *gameInfo,
				Message:  http.StatusText(http.StatusCreated),
			}
			if err := ws.WriteJSON(resp); err != nil {
				c.Logger().Error(err)
			}
			return err

		case CmdJoinGame:
			c.Logger().Print("Joining the game game...")
			if req.GameInfo.GameId == "" {
				errResp := models.Response{
					Code:  http.StatusBadRequest,
					Error: "No game id provided",
				}
				if err := ws.WriteJSON(errResp); err != nil {
					c.Logger().Error(err)
				}
				return err
			}

			gameInfo, err := JoinGame(req.GameInfo.GameId, c)
			if err != nil {
				errResp := models.Response{
					Code:  http.StatusInternalServerError,
					Error: err.Error(),
				}
				if err := ws.WriteJSON(errResp); err != nil {
					c.Logger().Error(err)
				}
				return err
			}

			resp := models.Response{
				Command:  CmdJoinGame,
				Code:     http.StatusOK,
				GameInfo: *gameInfo,
				Message:  http.StatusText(http.StatusOK),
			}
			if err := ws.WriteJSON(resp); err != nil {
				c.Logger().Error(err)
				return err
			}

		case CmdNewMove:
			c.Logger().Print("Making a move...")
			if req.GameInfo.GameId != "" && req.GameInfo.State != "" {
				errResp := models.Response{
					Code:  http.StatusBadRequest,
					Error: "No game state provided",
				}
				if err := ws.WriteJSON(errResp); err != nil {
					c.Logger().Error(err)
				}
				return err
			}

			gameInfo, err := NewMove(req.GameInfo, c)
			if err != nil {
				errResp := models.Response{
					Code:  http.StatusInternalServerError,
					Error: err.Error(),
				}
				if err := ws.WriteJSON(errResp); err != nil {
					c.Logger().Error(err)
				}
				return err
			}

			resp := models.Response{
				Command:  CmdNewMove,
				Code:     http.StatusOK,
				GameInfo: *gameInfo,
				Message:  http.StatusText(http.StatusOK),
			}
			if err := ws.WriteJSON(resp); err != nil {
				c.Logger().Error(err)
				return err
			}
		}
	}
}

func GenerateNewGame(c echo.Context) (*models.Game, error) {
	// Generate user id and game id
	userId := uuid.NewV4().String()
	c.Logger().Printf("user_id: %s\n", userId)
	gameId := uuid.NewV4().String()
	c.Logger().Printf("game_id: %s\n", gameId)

	repo := c.Get("GAME_REPO").(*repositories.GameRepo)
	g, err := repo.Create(gameId, userId)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func JoinGame(gameId string, c echo.Context) (*models.Game, error) {
	repo := c.Get("GAME_REPO").(*repositories.GameRepo)
	curGame, err := repo.FindByGameID(gameId)
	if err != nil {
		return nil, err
	}

	// Generate new user id
	userId := uuid.NewV4().String()
	c.Logger().Printf("user_id: %s\n", userId)

	// TODO: this fixes the bug with attempt to connect to ongoing game but for now keep it
	// if curGame.SecondUserId != "" {
	// 	err := errors.New("can't join game")
	// 	return nil, err
	// }

	// Register new user as second user
	curGame.SecondUserId = userId
	// Update game with new user
	err = repo.Update(curGame)
	if err != nil {
		return nil, err
	}

	return curGame, nil
}

func NewMove(game models.Game, c echo.Context) (*models.Game, error) {
	repo := c.Get("GAME_REPO").(*repositories.GameRepo)
	curGame, err := repo.FindByGameID(game.GameId)
	if err != nil {
		return nil, err
	}
	curGame.State = game.State
	curGame.LastMoveUserId = game.LastMoveUserId
	// Update game with new state
	err = repo.Update(curGame)
	if err != nil {
		return nil, err
	}

	return curGame, nil
}
