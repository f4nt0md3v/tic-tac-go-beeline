package ctx

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models/data"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models/game"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan data.Response
	Logger     *zap.SugaredLogger
}

func NewPool(logger *zap.SugaredLogger) *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan data.Response),
		Logger:     logger,
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			p.Logger.Info("Size of Connection Pool: ", len(p.Clients))
			for client := range p.Clients {
				_ = client.Conn.WriteJSON(data.Response{Message: "New User Joined..."})
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client)
			p.Logger.Info("Size of Connection Pool: ", len(p.Clients))
			for client := range p.Clients {
				_ = client.Conn.WriteJSON(data.Response{Message: "User Disconnected..."})
			}
			break
		case message := <-p.Broadcast:
			p.Logger.Info("Broadcast a message to all clients in Pool of connections...")
			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					p.Logger.Error(err)
					return
				}
			}
		}
	}
}

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Pool   *Pool
	AppCtx *AppContext
}

func NewClient(conn *websocket.Conn, ctx *AppContext) *Client {
	return &Client{
		Conn:   conn,
		Pool:   ctx.Pool,
		AppCtx: ctx,
	}
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			c.AppCtx.Logger.Error(err)
			return
		}

		var req data.Request
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			c.AppCtx.Logger.Error(err)
			return
		}
		c.AppCtx.Logger.Infof("Message Received: %+v\n", req)

		resp := ProcessRequest(&req, c.AppCtx)

		switch resp.MessageType {
		case data.Single:
			_ = c.Conn.WriteJSON(resp)
		case data.Broadcast:
			c.Pool.Broadcast <- *resp
		}
	}
}

const (
	CmdGenerateNewGame = "GENERATE_NEW_GAME"
	CmdJoinGame        = "JOIN_GAME"
	CmdNewMove         = "NEW_MOVE"
)

func ProcessRequest(req *data.Request, appCtx *AppContext) *data.Response {
	switch req.Command {
	case CmdGenerateNewGame:
		appCtx.Logger.Info("Generating a new game...")
		gameInfo, err := GenerateNewGame(appCtx)
		if err != nil {
			return &data.Response{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			}
		}

		resp := data.Response{
			Code:        http.StatusCreated,
			Command:     CmdGenerateNewGame,
			GameInfo:    gameInfo,
			Message:     http.StatusText(http.StatusCreated),
			MessageType: data.Single,
		}
		return &resp

	case CmdJoinGame:
		appCtx.Logger.Info("Joining the game game...")
		if req.GameInfo.GameId == "" {
			errResp := &data.Response{
				Code:  http.StatusBadRequest,
				Error: "No game id provided",
			}
			return errResp
		}

		gameInfo, err := JoinGame(req.GameInfo.GameId, appCtx)
		if err != nil {
			errResp := &data.Response{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			}
			return errResp
		}

		resp := &data.Response{
			Command:     CmdJoinGame,
			Code:        http.StatusOK,
			GameInfo:    gameInfo,
			Message:     http.StatusText(http.StatusOK),
			MessageType: data.Broadcast,
		}
		return resp

	case CmdNewMove:
		appCtx.Logger.Info("Making a move...")
		if req.GameInfo.GameId == "" && req.GameInfo.State == "" {
			errResp := &data.Response{
				Code:  http.StatusBadRequest,
				Error: "No game info provided",
			}
			return errResp
		}

		gameInfo, err := NewMove(*req.GameInfo, appCtx)
		if err != nil {
			errResp := &data.Response{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			}
			return errResp
		}

		resp := &data.Response{
			Command:     CmdNewMove,
			Code:        http.StatusOK,
			GameInfo:    gameInfo,
			Message:     http.StatusText(http.StatusOK),
			MessageType: data.Broadcast,
		}

		return resp
	}
	return nil
}

func GenerateNewGame(ctx *AppContext) (*game.Game, error) {
	// Generate user id and game id
	userId := uuid.NewV4().String()
	ctx.Logger.Infof("user_id: %s\n", userId)
	gameId := uuid.NewV4().String()
	ctx.Logger.Infof("game_id: %s\n", gameId)

	g, err := ctx.GameRepo.Create(gameId, userId)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func JoinGame(gameId string, ctx *AppContext) (*game.Game, error) {
	curGame, err := ctx.GameRepo.FindByGameID(gameId)
	if err != nil {
		return nil, err
	}

	// Generate new user id
	userId := uuid.NewV4().String()
	ctx.Logger.Infof("user_id: %s\n", userId)

	// Register new user as second user
	curGame.SecondUserId = userId
	// Update game with new user
	err = ctx.GameRepo.Update(curGame)
	if err != nil {
		return nil, err
	}

	return curGame, nil
}

func NewMove(game game.Game, ctx *AppContext) (*game.Game, error) {
	curGame, err := ctx.GameRepo.FindByGameID(game.GameId)
	if err != nil {
		return nil, err
	}
	curGame.State = game.State
	curGame.LastMoveUserId = game.LastMoveUserId
	// Update game with new state
	err = ctx.GameRepo.Update(curGame)
	if err != nil {
		return nil, err
	}

	return curGame, nil
}
