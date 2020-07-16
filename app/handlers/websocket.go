package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/models/message"
)

var (
	upg = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

func WebsocketHandler(c echo.Context) error {
	ws, err := upg.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := ws.Close(); err != nil {
			c.Logger().Error(err)
		}
	}()

	for {
		var readMsg message.Message
		// receive a message using the codec
		if err := ws.ReadJSON(&readMsg); err != nil {
			c.Logger().Error(err)
			break
		}
		c.Logger().Debug("Received message:", readMsg.Message)

		// send a response
		sendM := message.Message{Message: "Hello, Client!"}
		if err := ws.WriteJSON(sendM); err != nil {
			log.Println(err)
			break
		}
	}
	return nil
}
