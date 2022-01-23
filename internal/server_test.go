package internal

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

const wsPath = "ws://0.0.0.0:4567"

func TestWSMove(t *testing.T) {
	sender := NewSender()
	exit := make(chan bool)
	go sender.Listen(exit)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WS(sender, w, r)
	}))
	defer server.Close()
	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()
	cursor := CursorPoint{
		X: 123,
		Y: 456,
	}
	if err != ws.WriteJSON(cursor) {
		t.Fatalf("%v", err)
	}

	receivedMessage := new(Message)
	err = ws.ReadJSON(&receivedMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if receivedMessage.Method != MethodMove {
		t.Errorf("expected method %s got %v", MethodMove, receivedMessage.Method)
	}
	if receivedMessage.CursorPoint.X != cursor.X {
		t.Errorf("expected X %d got %v", cursor.X, receivedMessage.CursorPoint.X)
	}
	if receivedMessage.CursorPoint.Y != cursor.Y {
		t.Errorf("expected Y %d got %v", cursor.Y, receivedMessage.CursorPoint.Y)
	}
}

func TestWSLeave(t *testing.T) {
	sender := NewSender()
	exit := make(chan bool)
	go sender.Listen(exit)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WS(sender, w, r)
	}))
	defer server.Close()
	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()
	exit <- true

	receivedMessage := new(Message)
	err = ws.ReadJSON(&receivedMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if receivedMessage.Method != MethodLeave {
		t.Errorf("expected method %s got %v", MethodLeave, receivedMessage.Method)
	}
	if receivedMessage.CursorPoint != nil {
		t.Errorf("expected CursorPoint %v got %v", nil, receivedMessage.CursorPoint.X)
	}
}

func getTemplatePath() string {
	path, _ := os.Getwd()
	return strings.ReplaceAll(path, "internal", "templates/index.html")
}

func TestIndexOk(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	writer := httptest.NewRecorder()
	Index(getTemplatePath(), wsPath, writer, request)
	res := writer.Result()
	defer res.Body.Close()
	if writer.Code != http.StatusOK {
		t.Errorf("expected status code 200 got %d", writer.Code)
	}
	data, _ := ioutil.ReadAll(res.Body)
	if !strings.Contains(string(data), strings.ReplaceAll(wsPath, "/", "\\/")) {
		t.Error("ws path not found in the template")
	}
}

func TestIndex404(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	writer := httptest.NewRecorder()
	Index(getTemplatePath(), wsPath, writer, request)
	res := writer.Result()
	defer res.Body.Close()
	if writer.Code != http.StatusNotFound {
		t.Errorf("expected status code 404 got %d", writer.Code)
	}
}

func TestIndex405(t *testing.T) {
	//nolint:noctx
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	writer := httptest.NewRecorder()
	Index(getTemplatePath(), wsPath, writer, request)
	res := writer.Result()
	defer res.Body.Close()
	if writer.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status code 405 got %d", writer.Code)
	}
}
