package pubsub

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

const (
	PUBLISH   = "publish"
	SUBSCRIBE = "subscribe"
)

type PubSub struct {
	Clients []Client
}

type Client struct {
	Id         string
	Connection *websocket.Conn
}

type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

// AddClient adds a client to PubSub instance
func (p *PubSub) AddClient(c Client) {
	p.Clients = append(p.Clients, c)
}

// HandleReceiveMessage fetches message and acts based on Action of the Message
func (p *PubSub) HandleReceiveMessage(c Client, messageType int, payload []byte) *PubSub {
	m := Message{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Printf("wrong message payload: %v", err)
	}

	switch m.Action {
	case PUBLISH:
		log.Println("need to handle publish action")
	case SUBSCRIBE:
		log.Println("need to handle subscribe action")
	default:
		log.Println("unknown action type")
	}

	return p
}
