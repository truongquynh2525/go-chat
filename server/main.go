package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsHandler struct {
	header   http.Header
	upgrader *websocket.Upgrader
}

func newWsHandler() *wsHandler {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &wsHandler{
		upgrader: upgrader,
	}
}

func (ws *wsHandler) Open(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	if ws == nil || ws.upgrader == nil {
		return nil, fmt.Errorf("missing ws client")
	}

	return ws.upgrader.Upgrade(w, r, ws.header)
}

func wsHandleFunc(w http.ResponseWriter, r *http.Request) {
	wsHandler := newWsHandler()

	wsConn, err := wsHandler.Open(w, r)
	if err != nil {
		fmt.Println("cannot make socket connection")
		return
	}
	defer wsConn.Close()

	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println("read message error: ", err)
			return
		}

		if string(msg) == "PING" {
			wsConn.WriteMessage(websocket.TextMessage, []byte("PONG"))
			fmt.Println("PONG")
		}
		if string(msg) == "PONG" {
			wsConn.WriteMessage(websocket.TextMessage, []byte("PING"))
			fmt.Println("PING")
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandleFunc)

	fmt.Println("server is opening on port 3000")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("start server error")
	}
}
