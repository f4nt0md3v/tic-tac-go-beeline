package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models/ctx"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models/data"
)

var (
	upg = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: comment out or remove next line on production
	}
)

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upg.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func WebsocketHandler(appCtx *ctx.AppContext, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrade(w, r)
	if err != nil {
		log.Println(err)
		_ = conn.WriteJSON(data.Response{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	c := ctx.NewClient(conn, appCtx)
	appCtx.Pool.Register <- c

	c.Read()
}
