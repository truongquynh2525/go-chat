package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	client *websocket.Conn
}

func (wc *WebsocketClient) Close() error {
	return wc.client.Close()
}

func (wc *WebsocketClient) OnRead() {
	for {
		_, msg, err := wc.client.ReadMessage()
		if err != nil {
			panic(err)
		}

		fmt.Printf("receive from server %v \n", string(msg))
	}
}

func (wc *WebsocketClient) OnSend(done chan struct{}) {
	ticket := time.NewTicker(5 * time.Second)
	defer ticket.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticket.C:
			if err := wc.client.WriteMessage(websocket.TextMessage, []byte("PING")); err != nil {
				fmt.Println("send error")
				return
			}
			fmt.Printf("send to server %v \n", "PING")
		}
	}
}

func NewSocketConnection(url url.URL) (*WebsocketClient, error) {
	client, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}

	return &WebsocketClient{
		client: client,
	}, nil
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", u.String())

	ws, err := NewSocketConnection(u)
	if err != nil {
		fmt.Println("failed to connect ", err)
		panic(err)
	}
	defer ws.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		ws.OnRead()
	}()

	ws.OnSend(done)
}
