package main

import (
	"fmt"

	"github.com/nickstrad/streamer/collection"
)

func main() {
	wsClient := collection.NewWebsocketClient()
	fmt.Println(wsClient.Hml)
}
