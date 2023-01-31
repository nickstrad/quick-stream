package main

import (
	"fmt"
)

func main() {
	wsClient := NewWebsocketClient()
	fmt.Printf("%v\n", wsClient.Hml)
}
