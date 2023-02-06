package websocket_client

import (
	"log"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	Conn *websocket.Conn
}

func NewWebsocketClient() (*WebsocketClient, error) {
	/*
	 * Connect To Websocket
	 */
	log.Println("Connecting to websocket endpoint....")
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.coincap.io/trades/binance", nil)
	if err != nil {
		return nil, err
	}

	return &WebsocketClient{Conn: c}, nil
}

func (w *WebsocketClient) Cleanup() error {
	return w.Conn.Close()
}
