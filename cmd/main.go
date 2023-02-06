package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	hmlHelper "github.com/nickstrad/streamer/cmd/hybrid_message_logger"
	"github.com/nickstrad/streamer/cmd/websocket_client"
	"github.com/rs/xid"
	"github.com/segmentio/kafka-go"
)

func main() {

	/*
	 * Setup Hybrid Message Logger
	 */

	log.Println("Creating hybrid message logger databases....")
	hml, err := hmlHelper.NewHybridMessageLogger("data/transient", "data/failed")
	if err != nil {
		log.Fatal(err)
	}
	c, err := websocket_client.NewWebsocketClient()

	if err != nil {
		log.Fatal(err)
	}

	/*
	 * 	Create thread to read from the websocket into input channel
	 */
	type messageMetadata struct {
		message []byte
		id      xid.ID
	}

	type trade struct {
		Exchange  string  `json:"exchange"`
		Base      string  `json:"base"`
		Quote     string  `json:"quote"`
		Direction string  `json:"direction"`
		Price     float64 `json:"price"`
		Volume    float64 `json:"volume"`
		Timestamp int64   `json:"timestamp"`
		PriceUsd  float64 `json:"priceUsd"`
	}
	input := make(chan messageMetadata)

	log.Println("Listening for messages....")
	go func(c *websocket_client.WebsocketClient, hml *hmlHelper.HybridMessageLogger) {
		for {
			_, message, err := c.Conn.ReadMessage()

			if err != nil {
				log.Println("Error reading message....")
				log.Println(err)
				break
			}
			// unmarshal the message
			var trade trade
			err = json.Unmarshal(message, &trade)
			log.Printf("message received: %v\n", trade)
			if err != nil {
				log.Println(err)
				continue
			}
			id := xid.New()
			log.Println("Adding event with id: " + id.String())
			err = hml.AddEvent(id, message)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("Adding event with id: " + id.String() + " to kafka input")

			input <- messageMetadata{
				message: message,
				id:      id,
			}
		}
	}(c, hml)

	/*
	 * 	Create Producer to Send events to MessageQueue alyer
	 */
	topic := "collection-event"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	go func(conn *kafka.Conn, hml *hmlHelper.HybridMessageLogger) {
		defer conn.Close()
		for metadata := range input {
			bytesWritten, err := conn.WriteMessages(
				kafka.Message{Value: metadata.message},
			)

			if err != nil {
				log.Println(err)
				continue
			}

			if bytesWritten == 0 {
				hml.MoveToFailed(metadata.id)
			}

			if bytesWritten > 0 {
				hml.RemoveEvent(metadata.id)
			}
		}
	}(conn, hml)

	defer func() {
		defer hml.Cleanup()
		defer close(input)
		defer c.Cleanup()
	}()
	select {}

	// if err != nil {
	// 	log.Fatal("failed to write messages:", err)
	// }

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
}
