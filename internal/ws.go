package internal

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	symbols         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sessionIDLength = 20
	bufferSize      = 1024
)

var buffer = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
}

func generateSessionID(length int) string {
	raw := make([]byte, length)
	for i := range raw {
		raw[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(raw)
}

func (c Client) answer(msg []byte) {
	location := new(CursorPoint)
	err := json.Unmarshal(msg, location)
	if err != nil {
		log.Printf("unmarshal error: %s", err.Error())
	}
	c.Sender.Send(&Message{
		SessionID:   c.SessionID,
		Method:      MethodMove,
		CursorPoint: location,
	})
}

func (c Client) errorAnswer(err error) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		log.Printf("client closed connection: %s", err.Error())
	}
	c.Sender.Delete(c)
	defer c.Connection.Close()
	c.Sender.Send(&Message{
		SessionID:   c.SessionID,
		Method:      MethodLeave,
		CursorPoint: nil,
	})
}

func (c Client) Read() {
	for {
		_, msg, err := c.Connection.ReadMessage()
		if err != nil {
			c.errorAnswer(err)

			break
		}
		c.answer(msg)
	}
}

func WS(sender SenderManager, writer http.ResponseWriter, request *http.Request) {
	conn, err := buffer.Upgrade(writer, request, nil)
	if err != nil {
		log.Printf("connection initialization failure %s", err.Error())

		return
	}
	defer conn.Close()
	client := Client{
		Connection:   conn,
		SessionID:    generateSessionID(sessionIDLength),
		InputChannel: make(chan *Message),
		Sender:       sender,
	}
	sender.Add(client)
	client.Read()
}
