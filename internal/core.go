package internal

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MethodMove  = "move"
	MethodLeave = "leave"
)

type SenderManager interface {
	Add(client Client)
	Delete(client Client)
	Listen(quit <-chan bool)
	Send(message *Message)
}

type Client struct {
	Connection   *websocket.Conn
	SessionID    string
	InputChannel chan *Message
	Sender       SenderManager
}

type CursorPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Message struct {
	SessionID string `json:"sessionId"`
	Method    string `json:"method"`
	*CursorPoint
}

type Sender struct {
	mu      *sync.Mutex
	clients map[Client]struct{}
	output  chan *Message
}

func (s *Sender) Add(client Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client] = struct{}{}
}

func (s *Sender) Delete(client Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client)
}

func (s *Sender) Send(message *Message) {
	s.output <- message
}

func (s *Sender) Listen(quit <-chan bool) {
	for {
		select {
		case message := <-s.output:
			for client := range s.clients {
				err := client.Connection.WriteJSON(message)
				if err != nil {
					log.Printf("sending message error: %s", err.Error())
				}
			}

		case <-quit:
			for client := range s.clients {
				go func(client Client) {
					data := &Message{
						SessionID:   client.SessionID,
						Method:      MethodLeave,
						CursorPoint: nil,
					}
					s.output <- data
				}(client)
			}
		}
	}
}

func NewSender() *Sender {
	return &Sender{
		mu:      new(sync.Mutex),
		clients: make(map[Client]struct{}),
		output:  make(chan *Message),
	}
}
