package collection

// const MEETUP_URI = "ws://stream.meetup.com/2/rsvps"

type WebsocketClient struct {
	Hml *HybridMessageLogger
}

func (client *WebsocketClient) NewWebsocketClient() *WebsocketClient {
	hml := NewHybridMessageLogger()
	return &WebsocketClient{Hml: hml}
}
